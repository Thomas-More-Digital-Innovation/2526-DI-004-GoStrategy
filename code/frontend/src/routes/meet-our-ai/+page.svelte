<script lang="ts">
    import {
        aiDossiers,
        tournamentBenchmarks,
    } from "$lib/data/ai-intel.data";
    import AIDossierCard from "./_components/AIDossierCard.svelte";
    import AIComparisonTable from "./_components/AIComparisonTable.svelte";
    import AITournamentBenchmark from "./_components/AITournamentBenchmark.svelte";
    import Tabs from "$lib/components/ui/Tabs.svelte";

    type ViewMode = "dossiers" | "matrix" | "benchmarks";
    let activeView = $state<ViewMode>("dossiers");

    const tabs: { id: ViewMode; label: string }[] = [
        { id: "dossiers", label: "AI Dossiers" },
        { id: "matrix", label: "Algorithm Matrix" },
        { id: "benchmarks", label: "Tournament Benchmarks" },
    ];
</script>

<svelte:head>
    <title>GoStrategy — Meet our AI</title>
</svelte:head>

<div class="space-y-8 animate-fade-in pb-12">
    <!-- Page Header -->
    <header
        class="flex flex-col md:flex-row md:items-center justify-between gap-4"
    >
        <div>
            <h1
                class="text-3xl font-extrabold text-white uppercase tracking-widest"
            >
                Meet our AI
            </h1>
            <p class="text-white/50 text-sm mt-1">
                Explore the algorithms, heuristic models, and decision
                architectures powering GoStrategy opponents.
            </p>
        </div>
    </header>

    <Tabs {tabs} bind:active={activeView}>
        {#snippet children(current)}
            {#if current === "dossiers"}
                <div class="grid grid-cols-1 gap-6">
                    {#each aiDossiers as dossier}
                        <AIDossierCard {dossier} />
                    {/each}
                </div>
            {:else if current === "matrix"}
                <div class="space-y-6">
                    <AIComparisonTable dossiers={aiDossiers} />
                </div>
            {:else if current === "benchmarks"}
                <div class="space-y-6">
                    <AITournamentBenchmark benchmarks={tournamentBenchmarks} />
                </div>
            {/if}
        {/snippet}
    </Tabs>
</div>
