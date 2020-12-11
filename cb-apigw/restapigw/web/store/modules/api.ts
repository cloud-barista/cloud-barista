// =================================================================
//  Defines the module for API
// =================================================================

import { createModule, mutation, action } from "vuex-class-component";

import { ApiGroup, deserializeGroupFromJSON } from "@/models";
import {
  getApiGroups,
  addApiGroup,
  removeApiGroup,
  getApiGroup,
  updateApiDefinition,
  addApiDefinition,
  removeApiDefinition,
  applyChanges
} from "@/api/api";

const VuexModule = createModule({
  namespaced: "api",
  strict: false,
  target: "nuxt"
});

export default class ApiModule extends VuexModule {
  // ------------------------------------------------
  // Fields
  // ------------------------------------------------

  public groups: Array<ApiGroup> = [];

  // ------------------------------------------------
  // Mutations
  // ------------------------------------------------

  @mutation
  private SET_APIGROUPS(val: Array<ApiGroup>) {
    this.groups = val.map(v => deserializeGroupFromJSON(v, true));
  }

  @mutation
  private SET_APIGROUP(val: ApiGroup) {
    this.groups = this.groups.map(grp =>
      grp.name === val.name ? deserializeGroupFromJSON(val, true) : grp
    );
  }

  // ------------------------------------------------
  // Actions
  // ------------------------------------------------

  @action
  public async GetApiGroups() {
    const res = await getApiGroups();
    this.SET_APIGROUPS(res);
  }

  @action
  public async GetApiGroup(name: String) {
    return getApiGroup(name)
      .then(async res => {
        await this.SET_APIGROUP(res);
      })
      .catch(err => {
        Promise.reject(err);
      });
  }

  @action
  public async AddApiGroup(group: ApiGroup) {
    return addApiGroup(group)
      .then(async () => {
        await this.GetApiGroups();
      })
      .catch(err => {
        return Promise.reject(err);
      });
  }

  @action
  public async RemoveApiGroup(name: String) {
    return removeApiGroup(name)
      .then(async () => {
        await this.GetApiGroups();
      })
      .catch(err => {
        return Promise.reject(err);
      });
  }

  @action
  public async AddApiDefinition(group: ApiGroup) {
    return addApiDefinition(group)
      .then(async () => {
        await this.GetApiGroup(group.name);
      })
      .catch(err => {
        return Promise.reject(err);
      });
  }

  @action
  public async UpdateApiDefinition(group: ApiGroup) {
    return updateApiDefinition(group)
      .then(async () => {
        await this.GetApiGroup(group.name);
      })
      .catch(err => {
        return Promise.reject(err);
      });
  }

  @action
  public async RemoveApiDefinition(group: ApiGroup) {
    return removeApiDefinition(group)
      .then(async () => {
        await this.GetApiGroup(group.name);
      })
      .catch(err => {
        return Promise.reject(err);
      });
  }

  @action
  public async ApplyChanges() {
    return applyChanges()
      .then(async () => {
        await this.GetApiGroups();
      })
      .catch(err => {
        return Promise.reject(err);
      });
  }
}
