// Copyright 2016 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hugolib

import (
	"fmt"
)

// PageCollections contains the page collections for a site.
type PageCollections struct {
	// Includes only pages of all types, and only pages in the current language.
	Pages Pages

	// Includes all pages in all languages, including the current one.
	// Inlcudes pages of all types.
	AllPages Pages

	// A convenience cache for the traditional index types, taxonomies, home page etc.
	// This is for the current language only.
	indexPages Pages

	// A convenience cache for the regular pages.
	// This is for the current language only.
	// TODO(bep) np consider exporting this
	regularPages Pages

	// Includes absolute all pages (of all types), including drafts etc.
	rawAllPages Pages
}

func (c *PageCollections) refreshPageCaches() {
	c.indexPages = c.findPagesByNodeTypeNotIn(PagePage, c.Pages)
	c.regularPages = c.findPagesByNodeTypeIn(PagePage, c.Pages)

	// TODO(bep) np remove eventually
	for _, n := range c.Pages {
		if n.PageType == pageUnknown {
			panic(fmt.Sprintf("Got unknown type %s", n.Title))
		}
	}
}

func newPageCollections() *PageCollections {
	return &PageCollections{}
}

func newPageCollectionsFromPages(pages Pages) *PageCollections {
	return &PageCollections{rawAllPages: pages}
}

// TODO(bep) np clean and remove finders

func (c *PageCollections) findPagesByNodeType(n PageType) Pages {
	return c.findPagesByNodeTypeIn(n, c.Pages)
}

func (c *PageCollections) getPage(n PageType, path ...string) *Page {
	pages := c.findPagesByNodeTypeIn(n, c.Pages)

	if len(pages) == 0 {
		return nil
	}

	if len(path) == 0 && len(pages) == 1 {
		return pages[0]
	}

	for _, p := range pages {
		match := false
		for i := 0; i < len(path); i++ {
			if len(p.sections) > i && path[i] == p.sections[i] {
				match = true
			} else {
				match = false
				break
			}
		}
		if match {
			return p
		}
	}

	return nil
}

func (c *PageCollections) findIndexNodesByNodeType(n PageType) Pages {
	return c.findPagesByNodeTypeIn(n, c.indexPages)
}

func (*PageCollections) findPagesByNodeTypeIn(n PageType, inPages Pages) Pages {
	var pages Pages
	for _, p := range inPages {
		if p.PageType == n {
			pages = append(pages, p)
		}
	}
	return pages
}

func (*PageCollections) findPagesByNodeTypeNotIn(n PageType, inPages Pages) Pages {
	var pages Pages
	for _, p := range inPages {
		if p.PageType != n {
			pages = append(pages, p)
		}
	}
	return pages
}

func (c *PageCollections) findAllPagesByNodeType(n PageType) Pages {
	return c.findPagesByNodeTypeIn(n, c.Pages)
}

func (c *PageCollections) findRawAllPagesByNodeType(n PageType) Pages {
	return c.findPagesByNodeTypeIn(n, c.rawAllPages)
}

func (c *PageCollections) addPage(page *Page) {
	c.rawAllPages = append(c.rawAllPages, page)
}

func (c *PageCollections) removePageByPath(path string) {
	if i := c.rawAllPages.FindPagePosByFilePath(path); i >= 0 {
		c.rawAllPages = append(c.rawAllPages[:i], c.rawAllPages[i+1:]...)
	}
}

func (c *PageCollections) removePage(page *Page) {
	if i := c.rawAllPages.FindPagePos(page); i >= 0 {
		c.rawAllPages = append(c.rawAllPages[:i], c.rawAllPages[i+1:]...)
	}
}

func (c *PageCollections) replacePage(page *Page) {
	// will find existing page that matches filepath and remove it
	c.removePage(page)
	c.addPage(page)
}