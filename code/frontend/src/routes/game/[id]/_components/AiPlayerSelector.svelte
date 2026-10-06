<script lang="ts">
    import { useGameSession } from "../_state/context";

    interface Props {
        selectedPlayer: number;
        onSelectPlayer?: (player: number) => void;
    }

    let { selectedPlayer, onSelectPlayer }: Props = $props();

    const session = useGameSession();
</script>

<button
    class="flex items-center gap-3 mt-2 cursor-pointer group"
    aria-label="Select AI player"
    onclick={() => onSelectPlayer?.(selectedPlayer === 0 ? 1 : 0)}
>
    <div
        class="flex items-center gap-2 px-3 py-1.5 rounded-lg border transition-all {selectedPlayer ===
        0
            ? 'border-red-500/80 bg-red-500/15'
            : 'border-white/10 bg-white/5 opacity-50 group-hover:opacity-100 hover:border-white/20'}"
    >
        <div class="w-2 h-2 rounded-full bg-red-500"></div>
        <span
            class="text-[10px] font-bold text-white uppercase tracking-wider"
        >
            {session.store.gameState?.player1Username || "AI Red"}
        </span>
    </div>
    <div class="text-[10px] font-bold text-white/20 uppercase tracking-widest">vs</div>
    <div
        class="flex items-center gap-2 px-3 py-1.5 rounded-lg border transition-all {selectedPlayer ===
        1
            ? 'border-blue-500/80 bg-blue-500/15'
            : 'border-white/10 bg-white/5 opacity-50 group-hover:opacity-100 hover:border-white/20'}"
    >
        <span
            class="text-[10px] font-bold text-white uppercase tracking-wider"
        >
            {session.store.gameState?.player2Username || "AI Blue"}
        </span>
        <div class="w-2 h-2 rounded-full bg-blue-500"></div>
    </div>
</button>
