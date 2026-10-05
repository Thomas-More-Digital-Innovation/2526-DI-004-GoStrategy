<script lang="ts">
    import type { AITournamentBenchmark } from "$lib/types/ai";

    let { benchmarks }: { benchmarks: AITournamentBenchmark[] } = $props();

    const fafo = $derived(benchmarks.find((b) => b.aiId === "fafo"));
    const fato = $derived(benchmarks.find((b) => b.aiId === "fato"));
</script>

<div class="border border-white/10 bg-surface-elevated/40 rounded-xl overflow-hidden space-y-0">
    <div class="p-5 border-b border-white/10 flex flex-col sm:flex-row sm:items-center justify-between gap-2">
        <div>
            <h2 class="text-base font-bold text-white uppercase tracking-wider">
                10,000 Match Tournament Benchmark
            </h2>
            <p class="text-xs text-white/50 mt-0.5">
                Head-to-head simulation results comparing baseline random moves (FAFO) against piece memory (FATO).
            </p>
        </div>
    </div>

    {#if fafo && fato}
        <div class="overflow-x-auto">
            <table class="w-full text-left text-xs">
                <thead>
                    <tr class="border-b border-white/10 bg-black/20 text-[11px] uppercase tracking-wider text-white/50">
                        <th class="py-3 px-4 font-semibold">Metric</th>
                        <th class="py-3 px-4 font-semibold">FAFO (Random)</th>
                        <th class="py-3 px-4 font-semibold">FATO (Piece Memory)</th>
                        <th class="py-3 px-4 font-semibold">Observed Delta</th>
                    </tr>
                </thead>
                <tbody class="divide-y divide-white/5 font-mono">
                    <tr class="hover:bg-white/5 transition-colors">
                        <td class="py-3 px-4 font-sans text-white/70">Sample Size</td>
                        <td class="py-3 px-4 text-white">{fafo.sampleSize.toLocaleString()} games</td>
                        <td class="py-3 px-4 text-white">{fato.sampleSize.toLocaleString()} games</td>
                        <td class="py-3 px-4 text-white/40">—</td>
                    </tr>
                    <tr class="hover:bg-white/5 transition-colors">
                        <td class="py-3 px-4 font-sans text-white/70">Total Execution Time</td>
                        <td class="py-3 px-4 text-white">{fafo.totalRuntimeSeconds}s</td>
                        <td class="py-3 px-4 text-white">{fato.totalRuntimeSeconds}s</td>
                        <td class="py-3 px-4 text-white/80">+3.4x compute</td>
                    </tr>
                    <tr class="hover:bg-white/5 transition-colors">
                        <td class="py-3 px-4 font-sans text-white/70">Average Rounds / Game</td>
                        <td class="py-3 px-4 text-white">{fafo.avgRoundsPerGame.toFixed(1)}</td>
                        <td class="py-3 px-4 text-white">{fato.avgRoundsPerGame.toFixed(1)}</td>
                        <td class="py-3 px-4 text-white/80">-29.3% rounds</td>
                    </tr>
                    <tr class="hover:bg-white/5 transition-colors">
                        <td class="py-3 px-4 font-sans text-white/70">Flag Captured</td>
                        <td class="py-3 px-4 text-white">{fafo.flagCaptureRate} ({fafo.flagCaptures.toLocaleString()})</td>
                        <td class="py-3 px-4 text-white">{fato.flagCaptureRate} ({fato.flagCaptures.toLocaleString()})</td>
                        <td class="py-3 px-4 text-white/80">+8.5%</td>
                    </tr>
                    <tr class="hover:bg-white/5 transition-colors">
                        <td class="py-3 px-4 font-sans text-white/70">No Movable Pieces (Attrition)</td>
                        <td class="py-3 px-4 text-white">{fafo.noMoveWinsRate} ({fafo.noMoveWins.toLocaleString()})</td>
                        <td class="py-3 px-4 text-white">{fato.noMoveWinsRate} ({fato.noMoveWins.toLocaleString()})</td>
                        <td class="py-3 px-4 text-white/80">+11.7%</td>
                    </tr>
                    <tr class="hover:bg-white/5 transition-colors">
                        <td class="py-3 px-4 font-sans text-white/70">Max Turn Cutoffs (Stalemate)</td>
                        <td class="py-3 px-4 text-white">{fafo.maxTurnCutoffsRate} ({fafo.maxTurnCutoffs.toLocaleString()})</td>
                        <td class="py-3 px-4 text-white">{fato.maxTurnCutoffsRate} ({fato.maxTurnCutoffs.toLocaleString()})</td>
                        <td class="py-3 px-4 text-white/80">-20.1%</td>
                    </tr>
                </tbody>
            </table>
        </div>
    {/if}

    <div class="p-4 bg-black/20 border-t border-white/5 text-xs text-white/50 flex flex-col sm:flex-row justify-between gap-2">
        <span>Source: <code>documents/files/ai-data/fafo.md</code> and <code>fato.md</code></span>
        <span>Heuristic, Minimax, and MCTS benchmarks will be added once 10k tournament runs finish.</span>
    </div>
</div>
