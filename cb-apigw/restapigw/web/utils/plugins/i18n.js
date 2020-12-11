import Vue from "vue";
import VueI18n from "vue-i18n";

import { appStore } from "@/store";

Vue.use(VueI18n);

// export default ({ app, store }) => {
export default ({ app }) => {
  // Set i18n instance on app, This way we can use it in middleware and pages asyncData/fetch, ...
  app.i18n = new VueI18n({
    locale: appStore.locale,
    fallabckLocale: "en",
    messages: {
      en: require("~/utils/locales/en.json"),
      ko: require("~/utils/locales/ko.json")
    }
  });

  appStore.setLanguages([
    { id: "en", title: "English", flagSrc: "/images/lang/us.png" },
    { id: "ko", title: "Korean", flagSrc: "/images/lang/kr.png" }
  ]);

  // TODO: 필요성 여부 검증 필요.
  // Locale 기반 Routing Path 조정
  app.i18n.path = link => {
    if (app.i18n.locale === app.i18n.fallabckLocale) {
      return `/${link}`;
    }
    return `/${app.i18n.locale}/${link}`;
  };
};
