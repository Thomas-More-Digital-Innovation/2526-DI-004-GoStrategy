<script lang="ts">
    import type { AIDossier } from "$lib/types/ai";

    let { dossiers }: { dossiers: AIDossier[] } = $props();

    const trainabilityBadgeMap: Record<string, string> = {
        None: "text-red-400 bg-red-400/10 border-red-500/20",
        Partial: "text-amber-400 bg-amber-400/10 border-amber-500/20",
        Full: "text-emerald-400 bg-emerald-400/10 border-emerald-500/20",
    };
</script>

<div
    class="rounded-2xl border border-white/10 bg-surface-elevated/20 overflow-hidden backdrop-blur-md"
>
    <div class="p-6 border-b border-white/10">
        <h2 class="text-xl font-black text-white uppercase tracking-tight">
            Algorithmic Intel Matrix
        </h2>
        <p class="text-xs text-white/50 mt-1 uppercase tracking-wider">
            Comparative analysis of GoStrategy engine behaviors and search
            heuristics
        </p>
    </div>

    <div class="overflow-x-auto">
        <table class="w-full text-left text-xs border-collapse">
            <thead>
                <tr
                    class="border-b border-white/10 bg-white/5 text-[11px] uppercase tracking-wider text-white/60"
                >
                    <th class="py-3.5 px-4 font-bold">Agent</th>
                    <th class="py-3.5 px-4 font-bold">Info Model</th>
                    <th class="py-3.5 px-4 font-bold">Core Strength</th>
                    <th class="py-3.5 px-4 font-bold">Vulnerability</th>
                    <th class="py-3.5 px-4 font-bold">Trainability</th>
                    <th class="py-3.5 px-4 font-bold">Avg Latency</th>
                </tr>
            </thead>
            <tbody class="divide-y divide-white/5">
                {#each dossiers as ai}
                    <tr class="hover:bg-white/5 transition-colors">
                        <td
                            class="py-4 px-4 font-bold text-white flex items-center gap-3"
                        >
                            <img
                                src={ai.image}
                                alt={ai.name}
                                class="w-9 h-9 rounded-lg object-cover bg-black/40 border border-white/10"
                            />
                            <div>
                                <div class="text-sm uppercase tracking-tight">
                                    {ai.name}
                                </div>
                                <div
                                    class="text-[10px] text-white/40 font-normal"
                                >
                                    {ai.tagline}
                                </div>
                            </div>
                        </td>
                        <td class="py-4 px-4 text-white/80">
                            <span
                                class="px-2 py-0.5 rounded border border-white/10 bg-white/5 text-[10px] font-mono"
                            >
                                {ai.infoType}
                            </span>
                        </td>
                        <td class="py-4 px-4 text-emerald-300 font-medium">
                            {ai.strengths[0]}
                        </td>
                        <td class="py-4 px-4 text-brand-secondary font-medium">
                            {ai.weaknesses[0]}
                        </td>
                        <td class="py-4 px-4">
                            <span
                                class="px-2.5 py-0.5 rounded-full border text-[10px] font-bold tracking-wider uppercase {trainabilityBadgeMap[
                                    ai.trainability
                                ] || 'text-white'}"
                            >
                                {ai.trainability}
                            </span>
                        </td>
                        <td class="py-4 px-4 font-mono text-white/70">
                            {ai.stats.avgMoveTime || "—"}
                        </td>
                    </tr>
                {/each}
            </tbody>
        </table>
    </div>
</div>
