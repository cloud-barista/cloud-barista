import { $axios } from "@/utils/axios";
import { ApiGroup } from "@/models";

const rootUrl = "/apis";

/**
 * Get all group datas
 */
export const getApiGroups = () =>
  $axios.$get(`${rootUrl}/`).then(res => {
    return res;
  });

/**
 * Get group
 * @param name group name
 */
export const getApiGroup = (name: String) =>
  $axios.$get(`${rootUrl}/group/${name}`);

/**
 * Add group
 * @param name Group Object
 */
export const addApiGroup = (group: ApiGroup) =>
  $axios.$post(`${rootUrl}/group`, group);

/**
 * Remove group
 * @param name group name
 */
export const removeApiGroup = (name: String) => {
  // return $axios.$delete(rootUrl + "/group", { data: { Name: name } });
  return $axios.$delete(`${rootUrl}/group/${name}`);
};

/**
 * Add Definition
 * @param group Group Object
 */
export const addApiDefinition = (group: ApiGroup) => {
  return $axios.$post(`${rootUrl}/group/${group.name}/definition`, group);
};

/**
 * Update Api Definition
 * @param group Group Object
 */
export const updateApiDefinition = (group: ApiGroup) => {
  return $axios.$put(`${rootUrl}/group/${group.name}/definition`, group);
};

/**
 * Remove Definition
 * @param group Group Object
 */
export const removeApiDefinition = (group: ApiGroup) => {
  return $axios.$delete(
    `${rootUrl}/group/${group.name}/definition/${group.definitions[0].name}`
  );
};

/**
 * Apply Changes
 */
export const applyChanges = () => {
  return $axios.$put(`${rootUrl}/`);
};
