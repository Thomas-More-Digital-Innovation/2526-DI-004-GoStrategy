<script lang="ts">
    import Button from "$lib/components/ui/Button.svelte";
    import { serverStore } from "$lib/state/server.svelte";
    import { Zap, Unplug, Search } from "@lucide/svelte";

    interface Props {
        isReconnecting: boolean;
        gameId: string;
        reconnectAttempts: number;
        maxReconnectAttempts: number;
        onRetry?: () => void;
        onReturnToMenu?: () => void;
    }

    let {
        isReconnecting,
        gameId,
        reconnectAttempts,
        maxReconnectAttempts,
        onRetry = () => window.location.reload(),
        onReturnToMenu = () => (window.location.href = "/"),
    }: Props = $props();
</script>

<div
    class="flex flex-col items-center justify-center min-h-[60vh] gap-6 max-w-md mx-auto text-center px-4"
>
    {#if isReconnecting}
        <div class="relative flex items-center justify-center w-20 h-20">
            <div
                class="absolute inset-0 rounded-full border-2 border-brand-accent/20 border-t-brand-accent animate-spin"
            ></div>
            <Zap class="size-7 text-brand-accent" />
        </div>
        <div class="space-y-1.5">
            <h2 class="text-xl font-bold text-white tracking-wide uppercase">
                Connection Lost
            </h2>
            <p class="text-white/60 text-sm">
                Reconnecting to game session <span
                    class="font-mono text-brand-accent font-semibold"
                    >{gameId}</span
                >...
            </p>
            <p class="text-white/40 text-xs font-mono">
                Attempt {reconnectAttempts} of {maxReconnectAttempts}
            </p>
        </div>
    {:else if !serverStore.isOnline}
        <div
            class="flex items-center justify-center w-16 h-16 rounded-xl bg-brand-secondary/10 border border-brand-secondary/30 text-brand-secondary"
        >
            <Unplug class="size-7" />
        </div>
        <div class="space-y-1.5">
            <h2 class="text-xl font-bold text-white tracking-wide uppercase">
                Offline
            </h2>
            <p class="text-white/60 text-sm">
                Could not connect to the game server.
            </p>
        </div>
        <div class="flex flex-col sm:flex-row gap-3 w-full mt-2">
            <Button variant="outline" class="w-full" onclick={onRetry}>
                Retry Connection
            </Button>
            <Button variant="secondary" class="w-full" onclick={onReturnToMenu}>
                Return to Menu
            </Button>
        </div>
    {:else}
        <div
            class="flex items-center justify-center w-16 h-16 rounded-xl bg-brand-accent/10 border border-brand-accent/30 text-brand-accent"
        >
            <Search class="size-7" />
        </div>
        <div class="space-y-1.5">
            <h2 class="text-xl font-bold text-white tracking-wide uppercase">
                Game Not Found
            </h2>
            <p class="text-white/60 text-sm">
                Game session <span class="font-mono text-brand-accent font-semibold">{gameId}</span> could not be found.
            </p>
        </div>
        <div class="flex flex-col sm:flex-row gap-3 w-full mt-2">
            <Button variant="secondary" class="w-full" onclick={onReturnToMenu}>
                Return to Menu
            </Button>
        </div>
    {/if}
</div>
