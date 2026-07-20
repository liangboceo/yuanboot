import type { Directive, DirectiveBinding } from "vue";

interface NoSpaceInputElement extends HTMLElement {
  __noSpaceCleanup__?: () => void;
}

const emitInput = (input: HTMLInputElement) => {
  input.dispatchEvent(new Event("input", { bubbles: true }));
};

const findInput = (el: HTMLElement): HTMLInputElement | null => {
  if (el instanceof HTMLInputElement) return el;
  return el.querySelector("input");
};

const bindNoSpace = (el: NoSpaceInputElement) => {
  const input = findInput(el);
  if (!input) return;

  const onInput = () => {
    const nextValue = input.value.replace(/\s+/g, "");
    if (nextValue !== input.value) {
      input.value = nextValue;
      emitInput(input);
    }
  };

  const onCompositionEnd = () => onInput();

  input.addEventListener("input", onInput);
  input.addEventListener("compositionend", onCompositionEnd);

  el.__noSpaceCleanup__ = () => {
    input.removeEventListener("input", onInput);
    input.removeEventListener("compositionend", onCompositionEnd);
  };
};

export const noSpace: Directive = {
  mounted(el: NoSpaceInputElement) {
    bindNoSpace(el);
  },
  updated(el: NoSpaceInputElement, binding: DirectiveBinding) {
    if (binding.value === binding.oldValue) return;
    el.__noSpaceCleanup__?.();
    bindNoSpace(el);
  },
  unmounted(el: NoSpaceInputElement) {
    el.__noSpaceCleanup__?.();
    delete el.__noSpaceCleanup__;
  }
};
