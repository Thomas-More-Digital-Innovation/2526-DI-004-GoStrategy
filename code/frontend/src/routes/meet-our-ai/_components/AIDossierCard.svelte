<script lang="ts">
    import type { AIDossier } from "$lib/types/ai";
    import Button from "$lib/components/ui/Button.svelte";
    import { goto } from "$app/navigation";

    let { dossier }: { dossier: AIDossier } = $props();

    let expanded = $state(false);

    const infoTypeColorMap: Record<string, string> = {
        None: "border-zinc-500/40 text-zinc-300 bg-zinc-500/10",
        "Piece Memory": "border-cyan-500/40 text-cyan-300 bg-cyan-500/10",
        "Expert Rules": "border-amber-500/40 text-amber-300 bg-amber-500/10",
        Deterministic: "border-indigo-500/40 text-indigo-300 bg-indigo-500/10",
        Statistical: "border-emerald-500/40 text-emerald-300 bg-emerald-500/10",
    };

    const trainabilityColorMap: Record<string, string> = {
        None: "border-red-500/30 text-red-300 bg-red-500/10",
        Partial: "border-amber-500/30 text-amber-300 bg-amber-500/10",
        Full: "border-emerald-500/30 text-emerald-300 bg-emerald-500/10",
    };
</script>

<div
    class="flex flex-col bg-surface-elevated/20 border border-white/10 rounded-2xl p-6 backdrop-blur-md relative overflow-hidden group"
>
    <div class="flex items-start gap-5">
        <div class="relative shrink-0">
            <img
                src={dossier.image}
                alt={dossier.name}
                class="w-20 h-20 md:w-24 md:h-24 rounded-2xl object-cover bg-black/50 border border-white/10 shadow-lg group-hover:scale-102 transition-transform duration-300"
            />
            {#if dossier.stats.avgMoveTime}
                <div
                    class="absolute -bottom-2 -right-2 px-2 py-0.5 rounded-md text-[10px] font-mono font-bold uppercase tracking-wider bg-black/80 border border-white/20 text-white/80"
                >
                    {dossier.stats.avgMoveTime}
                </div>
            {/if}
        </div>

        <div class="flex-1 min-w-0">
            <div class="flex flex-wrap items-center gap-2 mb-1">
                <h3
                    class="text-2xl font-black text-white uppercase tracking-tight"
                >
                    {dossier.name}
                </h3>
                <span
                    class="px-2.5 py-0.5 text-[10px] font-semibold tracking-wider uppercase rounded-full border {infoTypeColorMap[
                        dossier.infoType
                    ] || 'border-white/20 text-white'}"
                >
                    {dossier.infoType}
                </span>
                <span
                    class="px-2.5 py-0.5 text-[10px] font-semibold tracking-wider uppercase rounded-full border {trainabilityColorMap[
                        dossier.trainability
                    ] || 'border-white/20 text-white'}"
                >
                    Trainable: {dossier.trainability}
                </span>
            </div>

            <p
                class="text-xs font-semibold text-brand-accent tracking-wide uppercase mb-2"
            >
                {dossier.tagline}
            </p>
            <p class="text-sm text-white/70 leading-relaxed">
                {dossier.description}
            </p>
        </div>
    </div>

    <!-- Metrics Matrix -->
    <div class="mt-6 pt-5 border-t border-white/10 grid grid-cols-3 gap-4">
        <div>
            <div
                class="flex justify-between items-center text-[11px] uppercase tracking-wider text-white/50 mb-1"
            >
                <span>Speed</span>
                <span class="font-mono text-white/90"
                    >{dossier.stats.speed}/5</span
                >
            </div>
            <div class="h-1.5 rounded-full bg-white/10 overflow-hidden">
                <div
                    class="h-full bg-brand-primary rounded-full transition-all duration-500"
                    style="width: {(dossier.stats.speed / 5) * 100}%"
                ></div>
            </div>
        </div>

        <div>
            <div
                class="flex justify-between items-center text-[11px] uppercase tracking-wider text-white/50 mb-1"
            >
                <span>Depth</span>
                <span class="font-mono text-white/90"
                    >{dossier.stats.strategicDepth}/5</span
                >
            </div>
            <div class="h-1.5 rounded-full bg-white/10 overflow-hidden">
                <div
                    class="h-full bg-indigo-400 rounded-full transition-all duration-500"
                    style="width: {(dossier.stats.strategicDepth / 5) * 100}%"
                ></div>
            </div>
        </div>

        <div>
            <div
                class="flex justify-between items-center text-[11px] uppercase tracking-wider text-white/50 mb-1"
            >
                <span>Adaptability</span>
                <span class="font-mono text-white/90"
                    >{dossier.stats.adaptability}/5</span
                >
            </div>
            <div class="h-1.5 rounded-full bg-white/10 overflow-hidden">
                <div
                    class="h-full bg-brand-accent rounded-full transition-all duration-500"
                    style="width: {(dossier.stats.adaptability / 5) * 100}%"
                ></div>
            </div>
        </div>
    </div>

    <!-- Strengths & Weaknesses -->
    <div class="mt-5 grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
        <div
            class="p-3.5 rounded-xl bg-black/20 border border-emerald-500/20 space-y-1.5"
        >
            <h4
                class="font-bold text-emerald-400 uppercase tracking-wider flex items-center gap-1.5 text-[11px]"
            >
                <span class="text-sm">✓</span> Strengths
            </h4>
            <ul class="space-y-1 text-white/70">
                {#each dossier.strengths as strength}
                    <li class="flex items-start gap-1.5 leading-snug">
                        <span class="text-emerald-400 shrink-0">•</span>
                        <span>{strength}</span>
                    </li>
                {/each}
            </ul>
        </div>

        <div
            class="p-3.5 rounded-xl bg-black/20 border border-brand-secondary/20 space-y-1.5"
        >
            <h4
                class="font-bold text-brand-secondary uppercase tracking-wider flex items-center gap-1.5 text-[11px]"
            >
                <span class="text-sm">✕</span> Weaknesses
            </h4>
            <ul class="space-y-1 text-white/70">
                {#each dossier.weaknesses as weakness}
                    <li class="flex items-start gap-1.5 leading-snug">
                        <span class="text-brand-secondary shrink-0">•</span>
                        <span>{weakness}</span>
                    </li>
                {/each}
            </ul>
        </div>
    </div>

    <!-- Deep-dive Collapsible -->
    {#if expanded}
        <div
            class="mt-5 p-4 rounded-xl bg-black/30 border border-white/10 space-y-3 animate-fade-in text-xs"
        >
            <div>
                <h5
                    class="text-[11px] font-bold text-brand-primary uppercase tracking-wider mb-1"
                >
                    Algorithmic Core
                </h5>
                <p class="text-white/80 leading-relaxed">
                    {dossier.concept}
                </p>
            </div>
            <div>
                <h5
                    class="text-[11px] font-bold text-brand-accent uppercase tracking-wider mb-1"
                >
                    Trainability & Optimization
                </h5>
                <p class="text-white/80 leading-relaxed">
                    {dossier.trainabilityDetails}
                </p>
            </div>
            {#if dossier.notes}
                <div class="pt-2 border-t border-white/5 text-white/50 italic">
                    {dossier.notes}
                </div>
            {/if}
        </div>
    {/if}

    <!-- Card Footer -->
    <div
        class="mt-6 pt-4 border-t border-white/10 flex items-center justify-between gap-3"
    >
        <button
            type="button"
            class="text-xs text-white/60 hover:text-white transition-colors cursor-pointer flex items-center gap-1"
            onclick={() => (expanded = !expanded)}
        >
            <span
                >{expanded
                    ? "Hide Details"
                    : "Inspect Algorithm Architecture"}</span
            >
            <span class="text-[10px]">{expanded ? "▲" : "▼"}</span>
        </button>

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
