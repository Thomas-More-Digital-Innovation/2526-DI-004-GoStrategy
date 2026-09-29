import type { AI } from "$lib/types/game";
import fafoImage from "$lib/assets/ai/fafo.webp";
import fatoImage from "$lib/assets/ai/fato.webp";
import heuristicImage from "$lib/assets/ai/heuristic.webp";
import minimaxImage from "$lib/assets/ai/minimax.webp";
import mctsImage from "$lib/assets/ai/mcts.webp";

export const AIs: AI[] = [
    {
        name: "FAFO",
        id: "fafo",
        description:
            "The Fuck Around & Find Out AI is a simple random-move AI.",
        image: fafoImage,
    },
    {
        name: "FATO",
        id: "fato",
        description:
            "The Find Around & Take Out AI uses piece memory and aggressive targeting.",
        image: fatoImage,
    },
    {
        name: "Heuristic",
        id: "heuristic",
        description:
            "Tactical Scout: Evaluates 1-ply board state, material balance, and territory advancement.",
        image: heuristicImage,
    },
    {
        name: "Minimax",
        id: "minimax",
        description:
            "Strategic Commander: Multi-ply lookahead with alpha-beta pruning and transposition tables.",
        image: minimaxImage,
    },
    {
        name: "MCTS",
        id: "mcts",
        description:
            "Monte Carlo Planner: Simulates rollouts to evaluate probabilistic outcomes.",
        image: mctsImage,
    },
];
