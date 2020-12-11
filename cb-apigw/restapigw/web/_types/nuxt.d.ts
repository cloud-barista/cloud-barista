import "@nuxt/vue-app/types";
import { Vue } from "nuxt-property-decorator";
import { VuetifyDialog } from "vuetify-dialog";
import "yamljs";

interface EventBus {
  listen(eventClass, handler): void;
  listenOnce(eventClass, handler): void;
  remove(eventClass, handler?): void;
  removeAll(): void;
  publish(event): void;
}

declare module "vue/types/vue" {
  interface Vue {
    // Default on Vue
    _uid: number;

    $dialog: VuetifyDialog;
    $eventBus: EventBus;
  }
}

declare module "vue/types/options" {
  interface ComponentOptions<V extends Vue> {
    auth?: boolean | string;
  }
}

declare global {
  interface Document {
    mozCancelFullscreen?: () => Promise<void>;
    webkitExitFullscreen?: () => Promise<void>;
    msExitFullscreen?: () => Promise<void>;
    mozFullScreenElement?: Element;
    msFullscreenElement?: Element;
    webkitFullscreenElement?: Element;
  }

  interface HTMLElement {
    mozRequestFullscreen?: () => Promise<void>;
    msRequestFullscreen?: () => Promise<void>;
    webkitRequestFullscreen?: () => Promise<void>;
  }

  interface Window {
    YAML: any;
  }
}
