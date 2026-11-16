/**
 * Encodes a 4x10 setup to a 40-character string of piece ranks.
 * Input rows use backend chars: 0=Flag, B=Bomb, 1=Spy, 2-9=Ranks, M=Marshal.
 */
export function encodeSetup(rows: string[]): string {
	// Ensure we have 4 rows of 10 chars, join into 40 chars
	let combined = '';
	for (let i = 0; i < 4; i++) {
		const row = rows[i] || '.'.repeat(10);
		combined += row.padEnd(10, '.').slice(0, 10);
	}
	return combined;
}

/**
 * Decodes a 40-character rank string to a 4x10 setup.
 * Returns rows using backend chars.
 */
export function decodeSetup(encoded: string): string[] {
	if (encoded.length !== 40) {
		// Fallback or handle error if needed, but we expect 40
		return Array(4).fill('.'.repeat(10));
	}
	const rows: string[] = [];
	for (let i = 0; i < 4; i++) {
		rows.push(encoded.slice(i * 10, (i + 1) * 10));
	}
	return rows;
}

/**
 * Flips a setup.
 */
export function flipSetup(encoded: string): string {
	const rows = decodeSetup(encoded);
	const flippedRows = rows.map((row) => row.split('').reverse().join('')).reverse();
	return encodeSetup(flippedRows);
}

export function formatBoardSetupForExport(setup: {
	name: string;
	description?: string;
	setup_data: string;
}) {
	return {
		name: setup.name,
		description: setup.description || '',
		setup_data: setup.setup_data,
		rows: decodeSetup(setup.setup_data)
	};
}

export function downloadJsonFile(filename: string, data: unknown): void {
	const json = JSON.stringify(data, null, 2);
	const blob = new Blob([json], { type: 'application/json' });
	const url = URL.createObjectURL(blob);
	const a = document.createElement('a');
	a.href = url;
	a.download = filename;
	document.body.appendChild(a);
	a.click();
	document.body.removeChild(a);
	URL.revokeObjectURL(url);
}

export function sanitizeFilename(name: string): string {
	const cleaned = name
		.toLowerCase()
		.trim()
		.replace(/[^a-z0-9_-]/g, '_');
	return cleaned || 'board_setup';
}

export function exportSetupAsJson(setup: {
	name: string;
	description?: string;
	setup_data: string;
}): void {
	const exported = formatBoardSetupForExport(setup);
	downloadJsonFile(`${sanitizeFilename(setup.name)}.json`, exported);
}
