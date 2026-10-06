<script lang="ts">
    import type { Snippet } from "svelte";
    import { fade, scale } from "svelte/transition";
    import { X } from "@lucide/svelte";

    interface Props {
        isOpen: boolean;
        title?: string;
        description?: string;
        onClose: () => void;
        maxWidth?: "sm" | "md" | "lg" | "xl" | "2xl" | "3xl" | "5xl";
        children: Snippet;
        actions?: Snippet;
        showCloseButton?: boolean;
        class?: string;
    }

    let {
        isOpen = $bindable(false),
        title = "",
        description = "",
        onClose,
        maxWidth = "lg",
        children,
        actions,
        showCloseButton = true,
        class: className = "",
    }: Props = $props();

    const widths = {
        sm: "max-w-sm",
        md: "max-w-md",
        lg: "max-w-lg",
        xl: "max-w-xl",
        "2xl": "max-w-2xl",
        "3xl": "max-w-3xl",
        "5xl": "max-w-5xl",
    };

    function handleKeydown(e: KeyboardEvent) {
        if (e.key === "Escape" && isOpen) {
            handleClose();
        }
    }

    function handleClose() {
        isOpen = false;
        onClose?.();
    }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if isOpen}
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
        transition:fade={{ duration: 180 }}
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/75 backdrop-blur-md"
        onclick={handleClose}
    >
        <div
            transition:scale={{ duration: 220, start: 0.96 }}
            class="relative w-full {widths[maxWidth]} flex flex-col max-h-[90vh] glass rounded-3xl overflow-hidden shadow-2xl border border-white/10 {className}"
            onclick={(e) => e.stopPropagation()}
        >
            {#if title || showCloseButton}
                <div
                    class="flex items-center justify-between border-b border-white/10 px-6 py-5 bg-white/5 shrink-0"
                >
                    <div class="space-y-0.5 pr-4">
                        {#if title}
                            <h2 class="text-lg font-bold text-white tracking-tight uppercase">
                                {title}
                            </h2>
                        {/if}
                        {#if description}
                            <p class="text-xs text-white/50">
                                {description}
                            </p>
                        {/if}
                    </div>

                    {#if showCloseButton}
                        <button
                            type="button"
                            aria-label="Close"
                            onclick={handleClose}
                            class="rounded-xl p-2 text-white/40 hover:text-white hover:bg-white/10 transition-colors cursor-pointer shrink-0"
                        >
                            <X class="size-4" />
                        </button>
                    {/if}
                </div>
            {/if}

            <div class="flex-1 overflow-y-auto p-6">
                {@render children()}
            </div>

            {#if actions}
                <div class="border-t border-white/10 px-6 py-4 bg-white/5 shrink-0">
                    {@render actions()}
                </div>
            {/if}
        </div>
    </div>
{/if}
