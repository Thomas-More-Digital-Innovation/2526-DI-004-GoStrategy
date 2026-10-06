<script lang="ts">
    import type { Snippet } from "svelte";
    import { Eye, EyeOff } from "@lucide/svelte";

    interface Props {
        type?: string;
        placeholder?: string;
        value?: string;
        label?: string;
        id?: string;
        class?: string;
        disabled?: boolean;
        sanitize?: "username" | "password" | "generic";
        error?: string;
        leadingIcon?: Snippet;
        showPasswordToggle?: boolean;
        focusColor?: "primary" | "secondary";
    }

    let {
        type = "text",
        placeholder = "",
        value = $bindable(""),
        label = "",
        id = Math.random().toString(36).substring(7),
        class: className = "",
        disabled = false,
        sanitize = undefined,
        error = "",
        leadingIcon,
        showPasswordToggle = false,
        focusColor = "primary",
    }: Props = $props();

    let passwordVisible = $state(false);

    const actualType = $derived(
        type === "password" ? (passwordVisible ? "text" : "password") : type,
    );

    function handleInput(e: Event) {
        const input = e.target as HTMLInputElement;
        let val = input.value;

        if (sanitize === "username") {
            val = val.replace(/[^a-zA-Z0-9_]/g, "");
        } else if (sanitize === "password") {
            val = val.replace(/(^ )|( $)|[^a-zA-Z0-9!@#$%^&*()_+=\-\. ]/g, "");
        } else if (sanitize === "generic") {
            val = val.replace(/[<>"'%;]/g, "");
        }

        if (val !== input.value) {
            input.value = val;
            value = val;
        }
    }
</script>

<div class="flex flex-col gap-1.5 {className}">
    {#if label}
        <label
            for={id}
            class="text-[10px] font-bold text-brand-accent uppercase tracking-widest ml-1"
        >
            {label}
        </label>
    {/if}

    <div class="relative group">
        {#if leadingIcon}
            <span
                class="absolute left-3.5 top-1/2 -translate-y-1/2 flex items-center pointer-events-none text-white/30 {focusColor === 'secondary' ? 'group-focus-within:text-brand-secondary' : 'group-focus-within:text-brand-primary'} transition-colors"
            >
                {@render leadingIcon()}
            </span>
        {/if}

        <input
            {id}
            type={actualType}
            {placeholder}
            {disabled}
            bind:value
            oninput={handleInput}
            class="w-full bg-white/5 border rounded-xl py-2.5 text-xs text-white placeholder:text-white/20 transition-all duration-200 focus:outline-none focus:ring-2 {leadingIcon ? 'pl-11' : 'px-4'} {showPasswordToggle && type === 'password' ? 'pr-11' : 'pr-4'} {error ? 'border-red-500/50 focus:ring-red-500/35 focus:border-red-500/35' : focusColor === 'secondary' ? 'border-white/10 focus:ring-brand-secondary/45 focus:border-brand-secondary/45' : 'border-white/10 focus:ring-brand-primary/45 focus:border-brand-primary/45'}"
        />

        {#if showPasswordToggle && type === "password"}
            <button
                type="button"
                aria-label={passwordVisible ? "Hide password" : "Show password"}
                onclick={() => (passwordVisible = !passwordVisible)}
                disabled={!value}
                class="absolute right-3 top-1/2 -translate-y-1/2 text-white/30 hover:text-white transition-colors cursor-pointer disabled:opacity-0 disabled:pointer-events-none p-1 rounded-md"
            >
                {#if passwordVisible}
                    <EyeOff class="size-4" />
                {:else}
                    <Eye class="size-4" />
                {/if}
            </button>
        {/if}
    </div>

    {#if error}
        <p class="text-[9px] text-red-400 font-bold uppercase tracking-wider ml-1 animate-fade-in">
            {error}
        </p>
    {/if}
</div>
