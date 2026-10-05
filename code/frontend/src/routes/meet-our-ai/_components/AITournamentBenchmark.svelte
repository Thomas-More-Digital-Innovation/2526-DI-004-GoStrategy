<script lang="ts">
    import type { AITournamentBenchmark } from "$lib/types/ai";

    let { benchmarks }: { benchmarks: AITournamentBenchmark[] } = $props();
</script>

<div class="space-y-6">
    <div
        class="rounded-2xl border border-white/10 bg-surface-elevated/20 p-6 backdrop-blur-md"
    >
        <div
            class="flex flex-col md:flex-row md:items-center justify-between gap-2 border-b border-white/10 pb-5 mb-6"
        >
            <div>
                <h2
                    class="text-xl font-black text-white uppercase tracking-tight"
                >
                    Empirical Tournament Benchmarks
                </h2>
                <p class="text-xs text-white/50 mt-1 uppercase tracking-wider">
                    Actual telemetry recorded from 10,000 automated tournament
                    games
                </p>
            </div>
            <span
                class="px-3 py-1 rounded-full text-[10px] font-mono font-bold tracking-wider uppercase border border-cyan-500/30 text-cyan-300 bg-cyan-500/10 w-fit"
            >
                Sample: 10,000 Matches / Agent
            </span>
        </div>

        <!-- Metric Comparison Cards -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            {#each benchmarks as bench}
                <div
                    class="p-5 rounded-xl bg-black/25 border border-white/10 space-y-4 flex flex-col justify-between"
                >
                    <div>
                        <div class="flex items-center justify-between mb-2">
                            <h3
                                class="text-lg font-black text-white uppercase tracking-tight"
                            >
                                {bench.aiName}
                            </h3>
                            <span
                                class="font-mono text-xs text-brand-primary font-bold"
                            >
                                {bench.totalRuntime} (10k games)
                            </span>
                        </div>
                        <p class="text-xs text-white/60 leading-relaxed mb-4">
                            {bench.notes}
                        </p>

                        <div class="grid grid-cols-2 gap-3 text-xs">
                            <div
                                class="p-3 rounded-lg bg-white/5 border border-white/5"
                            >
                                <div
                                    class="text-[10px] text-white/40 uppercase tracking-wider mb-1"
                                >
                                    Flag Capture
                                </div>
                                <div
                                    class="text-base font-black text-emerald-400 font-mono"
                                >
                                    {bench.flagCaptureRate}
                                </div>
                            </div>

                            <div
                                class="p-3 rounded-lg bg-white/5 border border-white/5"
                            >
                                <div
                                    class="text-[10px] text-white/40 uppercase tracking-wider mb-1"
                                >
                                    No Movable Pieces
                                </div>
                                <div
                                    class="text-base font-black text-cyan-300 font-mono"
                                >
                                    {bench.noMoveWinsRate}
                                </div>
                            </div>

                            <div
                                class="p-3 rounded-lg bg-white/5 border border-white/5"
                            >
                                <div
                                    class="text-[10px] text-white/40 uppercase tracking-wider mb-1"
                                >
                                    Max Turn Cutoffs
                                </div>
                                <div
                                    class="text-base font-black text-amber-300 font-mono"
                                >
                                    {bench.maxTurnCutoffsRate}
                                </div>
                            </div>

                            <div
                                class="p-3 rounded-lg bg-white/5 border border-white/5"
                            >
                                <div
                                    class="text-[10px] text-white/40 uppercase tracking-wider mb-1"
                                >
                                    Avg Rounds / Game
                                </div>
                                <div
                                    class="text-base font-black text-white font-mono"
                                >
                                    {bench.avgRoundsPerGame}
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    </div>

    <!-- Pending Benchmarks Card -->
    <div
        class="p-5 rounded-2xl border border-white/10 bg-surface-elevated/10 backdrop-blur-md flex flex-col md:flex-row md:items-center justify-between gap-4"
    >
        <div class="text-xs font-bold text-white uppercase tracking-wider">
            Heuristic, Minimax, MCTS & Mixed Benchmarks
        </div>
        <span
            class="px-3 py-1 rounded-full text-[10px] font-mono font-bold tracking-wider uppercase border border-amber-500/30 text-amber-300 bg-amber-500/10 w-fit shrink-0"
        >
            Awaiting Tournament
        </span>
    </div>
</div>
