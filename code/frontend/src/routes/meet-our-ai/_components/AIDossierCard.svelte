<script lang="ts">
    import type { AIDossier } from "$lib/types/ai";
    import Button from "$lib/components/ui/Button.svelte";
    import { goto } from "$app/navigation";

    let { dossier }: { dossier: AIDossier } = $props();

    let expanded = $state(false);
</script>

<div
    class="border border-white/10 bg-surface-elevated/40 rounded-xl p-5 space-y-4"
>
    <div class="flex flex-col sm:flex-row sm:items-start justify-between gap-4">
        <div class="flex items-start gap-4">
            <img
                src={dossier.image}
                alt={dossier.name}
                class="w-16 h-16 rounded-lg object-cover bg-black border border-white/10 shrink-0"
            />
            <div class="space-y-1">
                <div class="flex items-baseline gap-2">
                    <h3
                        class="text-xl font-bold text-white uppercase tracking-tight"
                    >
                        {dossier.name}
                    </h3>
                    <span class="text-xs text-white/50 font-mono">
                        {dossier.category}
                    </span>
                </div>
                <p class="text-sm text-white/80 max-w-2xl leading-relaxed">
                    {dossier.description}
                </p>
            </div>
        </div>

        <div class="sm:shrink-0">
            <Button
                onclick={() =>
                    goto(
                        `/select-ai?mode=human_vs_ai&ai=${encodeURIComponent(dossier.name)}`,
                    )}
                variant="primary"
                size="sm"
            >
                Challenge {dossier.name}
            </Button>
        </div>
    </div>

    <!-- Trade-offs Grid -->
    <div
        class="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs pt-3 border-t border-white/5"
    >
        <div
            class="p-3.5 rounded-lg bg-emerald-500/5 border border-emerald-500/20 space-y-1.5"
        >
            <div
                class="text-[11px] font-bold uppercase tracking-wider text-emerald-400 flex items-center gap-1.5"
            >
                <span class="text-xs">✓</span>
                <span>Strengths</span>
            </div>
            <ul class="space-y-1 text-white/80">
                {#each dossier.strengths as strength}
                    <li class="flex items-start gap-2">
                        <span class="text-emerald-400/80 shrink-0">•</span>
                        <span>{strength}</span>
                    </li>
                {/each}
            </ul>
        </div>

        <div
            class="p-3.5 rounded-lg bg-brand-secondary/5 border border-brand-secondary/20 space-y-1.5"
        >
            <div
                class="text-[11px] font-bold uppercase tracking-wider text-brand-secondary flex items-center gap-1.5"
            >
                <span class="text-xs">✕</span>
                <span>Limitations</span>
            </div>
            <ul class="space-y-1 text-white/80">
                {#each dossier.weaknesses as weakness}
                    <li class="flex items-start gap-2">
                        <span class="text-brand-secondary/80 shrink-0">•</span>
                        <span>{weakness}</span>
                    </li>
                {/each}
            </ul>
        </div>
    </div>

    <!-- Technical Details Expander -->
    {#if expanded}
        <div
            class="pt-3 border-t border-white/5 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 text-xs"
        >
            <div class="p-3 bg-black/20 rounded-lg border border-white/5">
                <div
                    class="text-white/40 text-[10px] uppercase tracking-wider mb-1"
                >
                    Algorithm
                </div>
                <div class="text-white/90 font-medium">{dossier.algorithm}</div>
            </div>
            <div class="p-3 bg-black/20 rounded-lg border border-white/5">
                <div
                    class="text-white/40 text-[10px] uppercase tracking-wider mb-1"
                >
                    State Model
                </div>
                <div class="text-white/90 font-medium">
                    {dossier.stateModel}
                </div>
            </div>
            <div class="p-3 bg-black/20 rounded-lg border border-white/5">
                <div
                    class="text-white/40 text-[10px] uppercase tracking-wider mb-1"
                >
                    Complexity
                </div>
                <div class="text-white/90 font-mono text-[11px]">
                    {dossier.complexity}
                </div>
            </div>
            <div class="p-3 bg-black/20 rounded-lg border border-white/5">
                <div
                    class="text-white/40 text-[10px] uppercase tracking-wider mb-1"
                >
                    Training Support
                </div>
                <div class="text-white/90 font-medium">
                    {dossier.trainingSupport}
                </div>
            </div>
        </div>
    {/if}

    <div
        class="pt-2 flex items-center justify-between text-xs text-white/50 border-t border-white/5"
    >
        <span>{dossier.notes || ""}</span>
        <button
            type="button"
            class="hover:text-white transition-colors cursor-pointer"
            onclick={() => (expanded = !expanded)}
        >
            {expanded ? "Hide specs" : "Show technical specs"}
        </button>
    </div>
</div>
