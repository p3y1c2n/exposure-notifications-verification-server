// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package admin

import (
	"context"
	"net/http"

	"github.com/google/exposure-notifications-verification-server/pkg/controller"
	"github.com/google/exposure-notifications-verification-server/pkg/database"
)

// HandleMobileAppsShow shows the configured mobile apps.
func (c *Controller) HandleMobileAppsShow() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		apps, err := c.db.ListActiveAppsWithRealm()
		if err != nil {
			controller.InternalError(w, r, c.h, err)
			return
		}

		c.renderShowMobileApps(ctx, w, apps)
	})
}

func (c *Controller) renderShowMobileApps(ctx context.Context, w http.ResponseWriter, apps []*database.ExtendedMobileApp) {
	m := controller.TemplateMapFromContext(ctx)
	m["apps"] = apps
	c.h.RenderHTML(w, "admin/mobileapps/show", m)
}
