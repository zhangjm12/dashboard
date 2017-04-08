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

/**
 * Controller for the replica set list view.
 *
 * @final
 */
export class ReleaseListController {
  /**
   * @param {!backendApi.ReleaseList} releaseList
   * @param {!angular.Resource} kdReleaseListResource
   * @ngInject
   */
  constructor(releaseList, kdReleaseListResource) {
    /** @export {!backendApi.ReleaseList} */
    this.releaseList = releaseList;

    /** @export {!angular.Resource} */
    this.releaseListResource = kdReleaseListResource;
  }

  /**
   * @return {boolean}
   * @export
   */
  shouldShowZeroState() {
    return this.releaseList.items.length === 0;
  }
}
