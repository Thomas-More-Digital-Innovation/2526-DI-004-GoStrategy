import { GameSocket } from '$lib/api/websocket';
import { games as gamesApi } from '$lib/api/client';
import { toastStore } from '$lib/state/toast.svelte';
import { serverStore } from '$lib/state/server.svelte';

export class GameConnectionManager {
	socket = new GameSocket();
	connected = $state(false);
	isReconnecting = $state(false);
	reconnectAttempts = $state(0);
	maxReconnectAttempts = 5;

	async connect(gameId: string, seatIndex: number) {
		try {
			await this.socket.connect(gameId, seatIndex);
			this.connected = true;
		} catch (e) {
			await serverStore.check();
			if (serverStore.isOnline) {
				toastStore.error('Game session not found or cleaned up.');
			} else {
				toastStore.handleApiMessage(e, 'Failed to connect to game server');
			}
			setTimeout(() => (window.location.href = '/'), 3000);
		}
	}

	async attemptReconnect(gameId: string, seatIndex: number) {
		if (this.isReconnecting) return;
		this.isReconnecting = true;
		this.reconnectAttempts = 0;

		while (this.reconnectAttempts < this.maxReconnectAttempts) {
			this.reconnectAttempts++;
			try {
				await this.socket.connect(gameId, seatIndex);
				this.connected = true;
				this.isReconnecting = false;
				this.reconnectAttempts = 0;
				toastStore.success('Reconnected to game session successfully.');
				return;
			} catch {
				const delay = Math.pow(2, this.reconnectAttempts) * 1000;
				await new Promise((resolve) => setTimeout(resolve, delay));
			}
		}

		this.isReconnecting = false;
		await serverStore.check();
		if (serverStore.isOnline) {
			toastStore.error('Failed to restore connection: Game session not found or cleaned up.');
		} else {
			toastStore.error('Failed to restore server connection after multiple attempts.');
		}
	}

	async abandonAndQuit(gameId: string) {
		try {
			this.disconnect();
			await gamesApi.abandon(gameId);
		} catch (e) {
			console.error('Failed to abandon session:', e);
		} finally {
			window.location.href = '/';
		}
	}

	disconnect() {
		this.socket.disconnect();
		this.connected = false;
	}
}
