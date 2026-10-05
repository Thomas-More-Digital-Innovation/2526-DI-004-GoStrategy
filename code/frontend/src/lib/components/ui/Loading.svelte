<script lang="ts">
    import { onMount } from "svelte";
    import fafo from "$lib/assets/ai/fafo.webp";
    import fato from "$lib/assets/ai/fato.webp";

    interface Props {
        title?: string;
        description?: string;
        subtitle?: string;
    }

    let {
        title = "Something is happening",
        description = "",
        subtitle = "Shhh... The computer is thinking",
    }: Props = $props();

    const images = [fafo, fato];
    let selectedImage = $state(fafo);

    onMount(() => {
        selectedImage = images[Math.floor(Math.random() * images.length)];
    });
</script>

<div
    class="flex flex-col items-center justify-center min-h-[50vh] gap-8 animate-in fade-in zoom-in-95 duration-700"
>
    <!-- Image Container -->
    <div class="relative">
        <div
            class="relative w-40 h-40 md:w-52 md:h-52 rounded-full overflow-hidden border border-white/10 bg-surface-elevated/60 shadow-xl"
        >
            <img
                src={selectedImage}
                alt={subtitle}
                class="w-full h-full object-contain animate-spin"
            />
        </div>
    </div>

    <!-- Text Content -->
    <div class="text-center space-y-3 max-w-sm">
        <div class="space-y-1">
            <h2
                class="text-3xl font-black text-white uppercase tracking-tighter"
            >
                {title}
            </h2>
            <p class="text-white/40 text-sm font-medium tracking-wide">
                {description}
            </p>
        </div>

        <div class="flex items-center justify-center gap-3 pt-4">
            <div
                class="h-px w-8 bg-linear-to-r from-transparent to-white/20"
            ></div>
            <span
                class="text-[10px] font-black text-brand-accent uppercase tracking-[0.3em] animate-pulse"
            >
                {subtitle}
            </span>
            <div
                class="h-px w-8 bg-linear-to-l from-transparent to-white/20"
            ></div>
        </div>
    </div>
</div>
