const toggleFullScreen = () => {
  const doc = window.document;
  const docEl = doc.documentElement;

  const requestFullScreen =
    docEl.requestFullscreen ||
    docEl.mozRequestFullscreen ||
    docEl.webkitRequestFullscreen ||
    docEl.msRequestFullscreen;
  const cancelFullScreen =
    doc.exitFullscreen ||
    doc.mozCancelFullscreen ||
    doc.webkitExitFullscreen ||
    doc.msExitFullscreen;

  if (
    !doc.fullscreenElement &&
    !doc.mozFullScreenElement &&
    !doc.webkitFullscreenElement &&
    !doc.msFullscreenElement
  ) {
    requestFullScreen.call(docEl);
  } else {
    cancelFullScreen.call(doc);
  }
};

const deserializeYaml = (yaml: string) => {
  if (yaml === "") return undefined;
  return window.YAML.parse(yaml);
};

const serializeYaml = (obj: any, depth: number = 10, indent: number = 2) => {
  if (!obj) return "";
  return window.YAML.stringify(obj, depth, indent);
};

class TimeParser {
  private millisecond = 1;
  private nano = this.millisecond / 1e6;
  private ms = this.millisecond;
  private second = this.millisecond * 1000;
  private s = this.second;
  private minute = this.s * 60;
  private m = this.minute;
  private hour = this.minute * 60;
  private h = this.hour;
  private duration = /(-?(?:\d+\.?\d*|\d*\.?\d+)(?:e[-+]?\d+)?)\s*([a-zµμ]*)/gi;

  public ToDuration(val: string) {
    if (typeof val === "number") val = String(val);

    let result = null;
    const str = val.replace(/(\d),(\d)/g, "$1$2");
    // console.log(`Inputed value: ${val}, processed val : ${str}`);
    result = str.replace(this.duration, (_, n, units) => {
      switch (units) {
        case "s":
          return ((parseFloat(n) * this.second) / this.nano).toString();
        case "m":
          return ((parseFloat(n) * this.minute) / this.nano).toString();
        case "h":
          return ((parseFloat(n) * this.hour) / this.nano).toString();
        default:
          return n;
      }
    });
    return Number(result);
  }

  public FromDuration(val: string, unit: string = "s") {
    let result = null;
    switch (unit) {
      case "s":
        result =
          ((parseFloat(val) * this.nano) / this.second).toString() + unit;
        break;
      case "m":
        result =
          ((parseFloat(val) * this.nano) / this.minute).toString() + unit;
        break;
      case "h":
        result = ((parseFloat(val) * this.nano) / this.hour).toString() + unit;
        break;
      default:
        result =
          ((parseFloat(val) * this.nano) / this.second).toString() + unit;
        break;
    }
    return result;
  }
}

const timeParser: TimeParser = new TimeParser();

export default {
  toggleFullScreen,
  timeParser,
  deserializeYaml,
  serializeYaml
};
