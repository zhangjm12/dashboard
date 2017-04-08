// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import stateConfig from './releaselist_stateconfig';
import filtersModule from 'common/filters/filters_module';
import componentsModule from 'common/components/components_module';
import namespaceModule from 'common/namespace/namespace_module';
import chromeModule from 'chrome/chrome_module';
import {releaseCardComponent} from './releasecard_component';
import {releaseCardListComponent} from './releasecardlist_component';
import releaseDetailModule from 'releasedetail/releasedetail_module';

/**
 * Angular module for the Replication Controller list view.
 *
 * The view shows Replication Controllers running in the cluster and allows to manage them.
 */
export default angular
    .module(
        'kubernetesDashboard.releaseList',
        [
          'ngMaterial',
          'ngResource',
          'ui.router',
          filtersModule.name,
          componentsModule.name,
          namespaceModule.name,
          chromeModule.name,
          releaseDetailModule.name,
        ])
    .config(stateConfig)
    .component('kdReleaseCardList', releaseCardListComponent)
    .component('kdReleaseCard', releaseCardComponent)
    .factory('kdReleaseListResource', releaseListResource);

/**
 * @param {!angular.$resource} $resource
 * @return {!angular.Resource}
 * @ngInject
 */
function releaseListResource($resource) {
  return $resource('api/v1/release/:namespace');
}
