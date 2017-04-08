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

import componentsModule from 'common/components/components_module';
import filtersModule from 'common/filters/filters_module';
import paginationModule from 'common/pagination/pagination_module';
import repositoryDetailModule from 'repositorydetail/repositorydetail_module';

import {RepositoryService} from './repository_service';
import {repositoryCardComponent} from './repositorycard_component';
import {repositoryCardListComponent} from './repositorycardlist_component';
import stateConfig from './repositorylist_stateconfig';


/**
 * Angular module for the Repository list view.
 *
 * The view shows Repository running in the cluster and allows to manage them.
 */
export default angular
    .module(
        'kubernetesDashboard.repositoryList',
        [
          'ngMaterial',
          'ngResource',
          'ui.router',
          filtersModule.name,
          componentsModule.name,
          repositoryDetailModule.name,
          paginationModule.name,
        ])
    .config(stateConfig)
    .component('kdRepositoryCardList', repositoryCardListComponent)
    .component('kdRepositoryCard', repositoryCardComponent)
    .service('kdRepositoryService', RepositoryService)
    .factory('kdRepositoryListResource', repositoryListResource);

/**
 * @param {!angular.$resource} $resource
 * @return {!angular.Resource}
 * @ngInject
 */
function repositoryListResource($resource) {
  return $resource('api/v1/repository');
}
