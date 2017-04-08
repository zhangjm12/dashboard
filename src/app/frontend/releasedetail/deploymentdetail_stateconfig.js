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

import {actionbarViewName, stateName as chromeStateName} from 'chrome/chrome_state';
import {breadcrumbsConfig} from 'common/components/breadcrumbs/breadcrumbs_service';
import {appendDetailParamsToUrl} from 'common/resource/resourcedetail';
import {stateName as deploymentList} from 'deploymentlist/deploymentlist_state';

import {ActionBarController} from './actionbar_controller';
import {DeploymentDetailController} from './deploymentdetail_controller';
import {stateName, stateUrl} from './deploymentdetail_state';

/**
 * Configures states for the deployment detail view.
 *
 * @param {!ui.router.$stateProvider} $stateProvider
 * @ngInject
 */
export default function stateConfig($stateProvider) {
  $stateProvider.state(stateName, {
    url: appendDetailParamsToUrl(stateUrl),
    parent: chromeStateName,
    resolve: {
      'deploymentDetailResource': getDeploymentDetailResource,
      'deploymentDetail': getDeploymentDetail,
    },
    data: {
      [breadcrumbsConfig]: {
        'label': '{{$stateParams.objectName}}',
        'parent': deploymentList,
      },
    },
    views: {
      '': {
        controller: DeploymentDetailController,
        controllerAs: 'ctrl',
        templateUrl: 'deploymentdetail/deploymentdetail.html',
      },
      [actionbarViewName]: {
        controller: ActionBarController,
        controllerAs: '$ctrl',
        templateUrl: 'deploymentdetail/actionbar.html',
      },
    },
  });
}

/**
 * @param {!./../common/resource/resourcedetail.StateParams} $stateParams
 * @param {!angular.$resource} $resource
 * @return {!angular.Resource<!backendApi.DeploymentDetail>}
 * @ngInject
 */
export function getDeploymentDetailResource($resource, $stateParams) {
  return $resource(`api/v1/deployment/${$stateParams.objectNamespace}/${$stateParams.objectName}`);
}

/**
 * @param {!angular.Resource<!backendApi.DeploymentDetail>} deploymentDetailResource
 * @return {!angular.$q.Promise}
 * @ngInject
 */
export function getDeploymentDetail(deploymentDetailResource) {
  return deploymentDetailResource.get().$promise;
}
