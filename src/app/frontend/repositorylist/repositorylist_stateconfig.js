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

import {actionbarViewName, stateName as chromeStateName} from 'chrome/state';
import {breadcrumbsConfig} from 'common/components/breadcrumbs/service';

import {ActionBarController} from './actionbar_controller';
import {RepositoryListController} from './repositorylist_controller';
import {stateName, stateUrl} from './repositorylist_state';

/**
 * Configures states for the service view.
 *
 * @param {!ui.router.$stateProvider} $stateProvider
 * @ngInject
 */
export default function stateConfig($stateProvider) {
  $stateProvider.state(stateName, {
    url: stateUrl,
    parent: chromeStateName,
    resolve: {
      'repositoryList': resolveRepositoryList,
    },
    data: {
      [breadcrumbsConfig]: {
        'label': i18n.MSG_BREADCRUMBS_REPOSITORIES_LABEL,
      },
    },
    views: {
      '': {
        controller: RepositoryListController,
        controllerAs: '$ctrl',
        templateUrl: 'repositorylist/repositorylist.html',
      },
      [actionbarViewName]: {
        controller: ActionBarController,
        controllerAs: '$ctrl',
        templateUrl: 'repositorylist/actionbar.html',
      },
    },
  });
}

/**
 * @param {!angular.$resource} $resource
 * @return {!angular.$q.Promise}
 * @ngInject
 */
export function resolveRepositoryList($resource) {
  /** @type {!angular.Resource<!backendApi.RepositoryList>} */
  let resource = $resource(`api/v1/repository`);
  return resource.get().$promise;
}

const i18n = {
  /** @type {string} @desc Label 'Repositories' that appears as a breadcrumbs on the action bar. */
  MSG_BREADCRUMBS_REPOSITORIES_LABEL: goog.getMsg('Repositories'),
};
