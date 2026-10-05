// TODO: source of truth should eventually come from a .json file in documents/files/ai-data/

import type { AIDossier, AITournamentBenchmark } from "$lib/types/ai";
import fafoImage from "$lib/assets/ai/fafo.webp";
import fatoImage from "$lib/assets/ai/fato.webp";
import heuristicImage from "$lib/assets/ai/heuristic.webp";
import minimaxImage from "$lib/assets/ai/minimax.webp";
import mctsImage from "$lib/assets/ai/mcts.webp";

export const aiDossiers: AIDossier[] = [
    {
        id: "fafo",
        name: "FAFO",
        tagline: "Pure Random Chaos Baseline",
        description:
            "The baseline agent selects moves uniformly at random from all legal actions without board awareness.",
        image: fafoImage,
        concept:
            "Evaluates all legal moves available in the current turn and picks one with uniform probability. Serves as our performance and regression baseline.",
        infoType: "None",
        trainability: "None",
        trainabilityDetails:
            "Random agents cannot be trained directly. They function as foundational sparring partners for training smarter models during RL warm-up phases.",
        strengths: [
            "Extremely fast with virtually zero computational overhead (~0.001ms/move).",
            "Completely unpredictable move vectors disrupt fixed opening patterns.",
            "Ideal baseline benchmark for measuring tactical intelligence gains.",
        ],
        weaknesses: [
            "No state evaluation, long-term memory, or tactical foresight.",
            "Frequently blunders high-value pieces and leaves flags undefended.",
            "Struggles against any agent with basic piece prioritization.",
        ],
        stats: {
            speed: 5,
            strategicDepth: 1,
            adaptability: 1,
            avgMoveTime: "< 0.01 ms",
        },
        notes: "Averaged 336.5 rounds per game in 100k simulation tournament; captures flag in ~62% of random wins.",
    },
    {
        id: "fato",
        name: "FATO",
        tagline: "Piece Memory Tactical Hunter",
        description:
            "Builds on random exploration with piece-memory caching and opportunistic targeting of revealed enemy pieces.",
        image: fatoImage,
        concept:
            "Maintains an internal memory map of observed enemy piece identities. When a known target is in range, it switches from random walk to tactical capture.",
        infoType: "Piece Memory",
        trainability: "None",
        trainabilityDetails:
            "Logic is deterministic heuristics combined with random walk exploration; weights are fixed but memory retention window is tunable.",
        strengths: [
            "Capitalizes immediately on revealed enemy vulnerabilities.",
            "Drastically reduces match length compared to pure random play.",
            "Low CPU overhead while punishing opponent scouts and exposed flags.",
        ],
        weaknesses: [
            "Cannot predict movement of unrevealed enemy units.",
            "Lacks deep positional understanding or multi-turn coordination.",
            "Can be baited into ambushes by higher-ranking hidden defenders.",
        ],
        stats: {
            speed: 5,
            strategicDepth: 2,
            adaptability: 2,
            avgMoveTime: "~0.05 ms",
        },
        notes: "Increases flag capture win rate to 70%+ in AI vs AI testing, reducing average match length to ~238 rounds.",
    },
    {
        id: "heuristic",
        name: "Heuristic",
        tagline: "Tactical Evaluation Scout",
        description:
            "Rule-based evaluator weighing material balance, board advancement, flag protection, and threat zones.",
        image: heuristicImage,
        concept:
            "Applies expert domain scoring across 1-ply candidate states: piece values, center control, threat proximity, and defensive screening.",
        infoType: "Expert Rules",
        trainability: "Partial",
        trainabilityDetails:
            "Evaluation weights can be optimized through genetic algorithms, hill-climbing, or parameter sweeps against game logs.",
        strengths: [
            "Calculates sharp, disciplined tactical decisions without lag.",
            "Encodes Stratego domain principles into piece positioning.",
            "Solid defensive structure protecting high-rank units and the flag.",
        ],
        weaknesses: [
            "Susceptible to designer bias and rigid positional assumptions.",
            "Struggles with imperfect information bluffing and deceptive scouts.",
            "Static evaluation cannot forecast deep multi-turn tactics.",
        ],
        stats: {
            speed: 4,
            strategicDepth: 3,
            adaptability: 2,
        },
        notes: "Serves as the primary tactical gatekeeper for testing deeper tree-search algorithms.",
    },
    {
        id: "minimax",
        name: "Minimax",
        tagline: "Strategic Alpha-Beta Commander",
        description:
            "Explores game trees with alpha-beta pruning and transposition tables, assuming optimal adversarial responses.",
        image: minimaxImage,
        concept:
            "Recursive adversarial search alternating between maximizing friendly payoff and minimizing opponent counter-moves up to a bounded depth limit.",
        infoType: "Deterministic",
        trainability: "Partial",
        trainabilityDetails:
            "The leaf evaluation function can be trained via self-play value nets, while the tree search mechanics remain algorithmic.",
        strengths: [
            "Theoretically optimal play under bounded lookahead horizons.",
            "Alpha-beta pruning aggressively eliminates unpromising move branches.",
            "Dominant in tactical endgames when piece positions are mostly exposed.",
        ],
        weaknesses: [
            "Imperfect information creates uncertainty that standard minimax struggles to model.",
            "Branching factor explodes without strict depth cutoffs.",
            "Vulnerable to bluffs when opponent ranks are unknown.",
        ],
        stats: {
            speed: 3,
            strategicDepth: 4,
            adaptability: 3,
        },
        notes: "Best deployed in deterministic phases or augmented with belief-state sampling.",
    },
    {
        id: "mcts",
        name: "MCTS",
        tagline: "Monte Carlo Probabilistic Planner",
        description:
            "Simulates rollout games per move, balancing exploration of unknown vectors and exploitation of proven wins.",
        image: mctsImage,
        concept:
            "Four-stage Monte Carlo Tree Search: Selection (UCT formula), Expansion, Simulation rollouts, and Backpropagation of win statistics.",
        infoType: "Statistical",
        trainability: "Full",
        trainabilityDetails:
            "Fully trainable with AlphaZero-style architecture: policy networks direct tree expansion while value networks replace random rollouts.",
        strengths: [
            "Excels in large, uncertain state spaces with hidden information.",
            "Dynamically converges towards optimal lines without hand-crafted heuristics.",
            "Foundational architecture for modern reinforcement learning engines.",
        ],
        weaknesses: [
            "High computational cost scaling with simulation budget.",
            "Rollout accuracy depends on playout policy quality.",
            "Requires careful time budgeting per turn in live interactive matches.",
        ],
        stats: {
            speed: 2,
            strategicDepth: 5,
            adaptability: 5,
        },
        notes: "Long-term champion architecture; scales smoothly with additional compute and trained neural policy networks.",
    },
];

// TODO: source of truth should eventually be loaded from documents/files/ai-data/*.json
export const tournamentBenchmarks: AITournamentBenchmark[] = [
    {
        aiId: "fafo",
        aiName: "FAFO",
        sampleSize: "10,000 matches",
        totalRuntime: "16.96 s",
        avgRoundsPerGame: "336.7",
        flagCaptureRate: "62.3%",
        noMoveWinsRate: "2.8%",
        maxTurnCutoffsRate: "34.9%",
        notes: "Baseline random mover; high rate of max-turn timeouts and accidental flag stumble captures.",
    },
    {
        aiId: "fato",
        aiName: "FATO",
        sampleSize: "10,000 matches",
        totalRuntime: "57.78 s",
        avgRoundsPerGame: "238.2",
        flagCaptureRate: "70.8%",
        noMoveWinsRate: "14.5%",
        maxTurnCutoffsRate: "14.8%",
        notes: "Piece memory increases flag captures by +8.5% and cuts match durations down by ~29.3%.",
    },
];
