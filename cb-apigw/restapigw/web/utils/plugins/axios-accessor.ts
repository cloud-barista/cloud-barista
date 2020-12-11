import { Plugin } from "@nuxt/types";
import { initializeAxios } from "@/utils/axios";

const accessor: Plugin = ({ $axios }) => {
  initializeAxios($axios);
};

export default accessor;
