<template>
  <div class="lined-textarea">
    <div
      v-if="!disabled"
      class="lined-textarea__lines"
      :style="{ 'padding-right': `${longestWidth}px` }"
    >
      <div ref="lines" class="lined-textarea__lines__inner">
        <p
          v-for="(line, index) in lines"
          :key="index"
          class="lined-textarea__lines__line"
          :class="{
            'lined-textarea__lines__line--invalid': invalidLines.includes(line)
          }"
          v-html="line"
        ></p>
      </div>
    </div>
    <textarea
      ref="textarea"
      v-model="content"
      :disabled="disabled"
      :placeholder="placeholder"
      class="lined-textarea__content"
      :class="{
        'lined-textarea__content--disabled': disabled,
        'lined-textarea__content--wrap': !nowrap,
        'lined-textarea__content--nowrap': nowrap
      }"
      :style="styles"
      @scroll="scrollLines"
      @input="onInput"
      @mousedown="detectResize"
    ></textarea>
    <div ref="helper" class="count-helper"></div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Prop, Watch } from "nuxt-property-decorator";

@Component({
  name: "cb-lined-textarea",
  components: {}
})
export default class CBLinedTextarea extends Vue {
  @Prop({ type: Boolean, default: false }) disabled!: boolean;
  @Prop({ type: Boolean, default: true }) nowrap!: boolean;
  @Prop({ type: String, default: "" }) placeholder!: string;
  @Prop({
    type: Object,
    default: () => {
      "300px";
    }
  })
  styles?: object;
  @Prop({ type: Array, default: () => [] }) errorLines!: Array<Number>;
  @Prop({ type: String, default: "" }) value!: string;
  @Prop({ type: Function, default: () => true }) validate!: Function;
  @Watch("longestWidth")
  onLongestWidthChanged(val: any, oldVal: any) {
    if (val !== oldVal) this.$nextTick(() => this.calculateCharactersPerLine());
  }
  @Watch("nowrap")
  onNoWrapChanged() {
    this.calculateCharactersPerLine();
  }
  @Watch("value")
  onValueChagned(val: any) {
    if (val !== this.content) this.content = val;
  }

  // ---------------------------------
  // Fields
  // ---------------------------------

  private content: string = "";
  private widthPerChar: number = 8;
  private numPerLine: number = 1;
  private previousWidth: number = 0;
  private scrolledLength: number = 0;

  // ---------------------------------
  // Getters/Setters
  // ---------------------------------

  private get invalidLines() {
    const lineNumbers: Array<any> = this.errorLines || [];
    this.content.split("\n").forEach((line, index) => {
      if (!this.validate(line)) lineNumbers.push(index + 1);
    });
    return lineNumbers;
  }

  private get lines() {
    if (this.content === "") return [1];
    const lineNumbers: Array<any> = [];
    let num = 1;
    // Number of lines extended. Seems to work with pre-wrap (has problem with dash)
    function getWrapTimes(sentence: any, width: any) {
      if (width <= 0) {
        // Protect against infinite loop
        console.warn(
          "Please set the min-width of textarea to allow at least one character per line."
        );
        return sentence.length + 1; // Seems browser would add one additional space
      }
      const words = sentence.split(" ");
      let currentLine = 1;
      let spaceLeft = width;
      words.forEach((word: any) => {
        const isWidth = spaceLeft === width;
        /* eslint-disable no-unmodified-loop-condition */
        while (isWidth && word.length >= spaceLeft) {
          ++currentLine;
          word = word.slice(width);
        }
        if (spaceLeft === width) {
          spaceLeft -= word.length;
          return;
        }
        if (word.length + 1 > spaceLeft) {
          ++currentLine;
          spaceLeft = width;
        }
        spaceLeft -= isWidth ? word.length : word.length + 1;
      });
      return spaceLeft === width ? currentLine - 1 : currentLine;
    }
    this.content.split("\n").forEach(line => {
      const wrapTimes = getWrapTimes(line, this.numPerLine) - 1;
      lineNumbers.push(num);
      for (let i = 0; i < wrapTimes; ++i) lineNumbers.push("&nbsp;");
      ++num;
    });
    return lineNumbers;
  }

  private get longestWidth() {
    for (let i = this.lines.length - 1; i >= 0; --i) {
      if (this.lines[i] === "&nbsp;") continue;
      return (this.lines[i] + "").length * this.widthPerChar + 10; // 10px base padding-right
    }
    return 5 * this.widthPerChar + 10;
  }

  // ---------------------------------
  // Methods
  // ---------------------------------

  private calculateCharactersPerLine() {
    if (this.nowrap) {
      this.numPerLine = Number.MAX_SAFE_INTEGER;
      return;
    }
    const textarea = this.$refs.textarea as HTMLTextAreaElement;
    const styles = getComputedStyle(textarea);
    const paddingLeft = parseFloat(styles.getPropertyValue("padding-left"));
    const paddingRight = parseFloat(styles.getPropertyValue("padding-right"));
    const fontSize = styles.getPropertyValue("font-size");
    const fontFamily = styles.getPropertyValue("font-family");
    const width = textarea.clientWidth - paddingLeft - paddingRight;
    const helper = this.$refs.helper as HTMLDivElement;
    helper.style.fontSize = fontSize;
    helper.style.fontFamily = fontFamily;
    helper.innerHTML = "";
    for (let num = 1; num < 999; ++num) {
      helper.innerHTML += "a";
      if (helper.getBoundingClientRect().width > width) {
        this.numPerLine = num - 1;
        break;
      }
    }
  }

  private detectResize() {
    const textarea = this.$refs.textarea as HTMLTextAreaElement;
    const { clientWidth: w1, clientHeight: h1 } = textarea;
    const detect = () => {
      const { clientWidth: w2, clientHeight: h2 } = textarea;
      if (w1 !== w2 || h1 !== h2) this.calculateCharactersPerLine();
    };
    const stop = () => {
      this.calculateCharactersPerLine();
      document.removeEventListener("mousemove", detect);
      document.removeEventListener("mouseup", stop);
    };
    document.addEventListener("mousemove", detect);
    document.addEventListener("mouseup", stop);
  }

  private onInput() {
    this.$emit("input", this.content);
    this.recalculate();
  }

  private recalculate() {
    const textarea = this.$refs.textarea as HTMLTextAreaElement;
    const width = textarea.clientWidth;
    if (width !== this.previousWidth) this.calculateCharactersPerLine();
    this.previousWidth = width;
  }

  private scrollLines(e: any) {
    this.scrolledLength = e.target.scrollTop;
    this.syncScroll();
  }

  private syncScroll() {
    if (!this.disabled) {
      (this.$refs.lines as HTMLDivElement).style.transform = `translateY(${-this
        .scrolledLength}px)`;
    }
  }

  // ---------------------------------
  // Lifecycle Events
  // ---------------------------------

  mounted() {
    this.content = this.value;
    this.syncScroll();
    this.calculateCharactersPerLine();
  }
}
</script>

<style lang="scss" scoped>
.lined-textarea {
  display: flex;
  font-size: 13px;
  line-height: 150%;
  font-family: Helvetica, monospace;
  background-color: white;
}
.lined-textarea__lines {
  background-color: #f0f0f0;
  border: 1px solid #d7e2ed;
  // border-radius: 10px 0 0 10px;
  border-right-width: 0;
  padding: 15px 10px 15px 15px;
  overflow: hidden;
  position: relative;
  color: black;
  font-size: 13px;
  line-height: 150%;
}

.lined-textarea__lines p {
  margin-bottom: 0px !important;
}

.lined-textarea__lines__inner {
  position: absolute;
}
.lined-textarea__lines__line {
  text-align: right;
}
.lined-textarea__lines__line--invalid {
  font-weight: bold;
  color: red;
}
.lined-textarea__content {
  border: 1px solid #d7e2ed;
  //  border-radius: 0 10px 10px 0;
  border-left-width: 0;
  margin: 0;
  line-height: inherit;
  font-family: monospace;
  padding: 15px;
  width: 100%;
  overflow: auto;
}
.lined-textarea__content--wrap {
  white-space: pre-wrap;
}
.lined-textarea__content--nowrap {
  white-space: pre;
}
@supports (-ms-ime-align: auto) {
  .lined-textarea__content--nowrap {
    white-space: nowrap;
  }
}
.lined-textarea__content--disabled {
  // border-radius: 10px;
  // border-left-width: 1px;
  background-color: lightgray;
}
.lined-textarea__content:focus {
  outline: none;
}
.count-helper {
  position: absolute;
  visibility: hidden;
  height: auto;
  width: auto;
}
</style>
