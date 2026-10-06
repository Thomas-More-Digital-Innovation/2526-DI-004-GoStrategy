// TODO: source of truth should eventually be loaded from documents/files/ai-data/*.json

import type { AIDossier, AITournamentBenchmark } from '$lib/types/ai';
import fafoImage from '$lib/assets/ai/fafo.webp';
import fatoImage from '$lib/assets/ai/fato.webp';
import heuristicImage from '$lib/assets/ai/heuristic.webp';
import minimaxImage from '$lib/assets/ai/minimax.webp';
import mctsImage from '$lib/assets/ai/mcts.webp';

export const aiDossiers: AIDossier[] = [
	{
		id: 'fafo',
		name: 'FAFO',
		category: 'Uniform Random',
		description:
			'Selects uniformly from all legal moves on each turn with zero board evaluation or state memory.',
		image: fafoImage,
		algorithm: 'Uniform random selection across legal action list.',
		stateModel: 'Stateless (no history, no piece tracking).',
		complexity: 'O(1) decision time.',
		trainingSupport: 'None',
		trainingDetails:
			'Cannot be trained. Functions as a baseline for regression tests and early-stage RL opponent pools.',
		strengths: [
			'Near-zero computational cost (<0.01ms per move).',
			'Unpredictable baseline useful for detecting degenerate play patterns.'
		],
		weaknesses: [
			'No tactical awareness or objective pursuit.',
			'High rate of moves into known enemy threats.'
		],
		notes: 'Averaged 336.7 rounds per game across 10,000 automated tournament matches.'
	},
	{
		id: 'fato',
		name: 'FATO',
		category: 'Heuristic with Piece Memory',
		description:
			'Extends random exploration with piece-memory caching, prioritizing immediate attacks on revealed enemy targets.',
		image: fatoImage,
		algorithm:
			'Random walk with deterministic override when revealed enemy pieces enter strike range.',
		stateModel: 'Local board state plus revealed enemy piece position and rank cache.',
		complexity: 'O(N) range check against tracked enemy pieces.',
		trainingSupport: 'None',
		trainingDetails:
			'Static priority rules. Memory retention window and targeting weights can be parameterized.',
		strengths: [
			'Directly exploits revealed enemy positions and exposed flags.',
			'Reduces match duration by ~29% compared to random play.'
		],
		weaknesses: [
			'Cannot infer positions or ranks of unrevealed pieces.',
			'Vulnerable to ambush by higher-ranking hidden defenders.'
		],
		notes:
			'Achieved 70.8% flag capture rate in 10k matches, lowering average game length to 238.2 rounds.'
	},
	{
		id: 'heuristic',
		name: 'Heuristic',
		category: '1-Ply Evaluation Engine',
		description:
			'Evaluates candidate 1-ply board states against expert-designed positional and material scoring rules.',
		image: heuristicImage,
		algorithm:
			'Linear evaluation function weighting material balance, territory advancement, and flag defense.',
		stateModel: 'Full 1-ply forward board state with rule-based features.',
		complexity: 'O(M) where M is the count of legal moves.',
		trainingSupport: 'Partial',
		trainingDetails:
			'Evaluation weights can be tuned via genetic algorithms, hill-climbing, or game log regression.',
		strengths: [
			'Low latency with consistent defensive piece structuring.',
			'Encodes standard Stratego principles without deep search overhead.'
		],
		weaknesses: [
			'No multi-ply forward lookahead.',
			'Vulnerable to tactical sacrifices and bluffing in hidden information states.'
		],
		notes: 'Benchmark evaluation currently pending.'
	},
	{
		id: 'minimax',
		name: 'Minimax',
		category: 'Alpha-Beta Adversarial Search',
		description:
			'Recursive game tree search assuming optimal play, pruned with alpha-beta bounds and transposition tables.',
		image: minimaxImage,
		algorithm: 'Depth-limited minimax search with alpha-beta cutoffs and transposition cache.',
		stateModel: 'Recursive board tree states with deterministic leaf evaluations.',
		complexity: 'O(b^(d/2)) best-case with optimal move ordering.',
		trainingSupport: 'Partial',
		trainingDetails:
			'Leaf evaluation function is trainable via self-play; search algorithm remains deterministic.',
		strengths: [
			'Optimal play within bounded lookahead horizons.',
			'Strong endgame tactical execution once pieces are revealed.'
		],
		weaknesses: [
			'Hidden information requires belief-state approximations.',
			'Branching factor limits search depth without aggressive pruning.'
		],
		notes: 'Benchmark evaluation currently pending.'
	},
	{
		id: 'mcts',
		name: 'MCTS',
		category: 'Monte Carlo Tree Search',
		description:
			'Simulates rollout games from current position, balancing exploration and exploitation to evaluate moves statistically.',
		image: mctsImage,
		algorithm: 'UCT tree search: Selection, Expansion, Simulation rollouts, and Backpropagation.',
		stateModel: 'Tree of visited board states with win/loss visit counts.',
		complexity: 'O(S × d) where S is simulation count and d is rollout depth.',
		trainingSupport: 'Full',
		trainingDetails:
			'Compatible with AlphaZero architectures: policy networks guide expansion, value networks replace rollouts.',
		strengths: [
			'Handles large state spaces and imperfect information without fixed heuristics.',
			'Scales decision quality directly with available compute budget.'
		],
		weaknesses: [
			'High CPU demand per decision compared to 1-ply heuristics.',
			'Requires rollout budget tuning for interactive latency constraints.'
		],
		notes: 'Benchmark evaluation currently pending.'
	}
];

export const tournamentBenchmarks: AITournamentBenchmark[] = [
	{
		aiId: 'fafo',
		aiName: 'FAFO',
		sampleSize: 10000,
		totalRuntimeSeconds: 16.96,
		avgRoundsPerGame: 336.67,
		flagCaptures: 5404,
		flagCaptureRate: '62.3%',
		noMoveWins: 243,
		noMoveWinsRate: '2.8%',
		maxTurnCutoffs: 3027,
		maxTurnCutoffsRate: '34.9%',
		summary: 'Baseline random mover; high rate of max-turn timeouts.'
	},
	{
		aiId: 'fato',
		aiName: 'FATO',
		sampleSize: 10000,
		totalRuntimeSeconds: 57.78,
		avgRoundsPerGame: 238.18,
		flagCaptures: 5328,
		flagCaptureRate: '70.8%',
		noMoveWins: 1090,
		noMoveWinsRate: '14.5%',
		maxTurnCutoffs: 1111,
		maxTurnCutoffsRate: '14.8%',
		summary: 'Piece memory reduces stall draws by 20% and accelerates flag captures.'
	}
];
