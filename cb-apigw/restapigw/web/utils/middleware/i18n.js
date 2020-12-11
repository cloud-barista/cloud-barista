import { appStore } from "@/store";

// export default async function({ app, store }) {
export default async function({ app }) {
  // Set Language
  app.i18n.locale = appStore.locale;
}
