export type AIInfoType =
    | "None"
    | "Piece Memory"
    | "Expert Rules"
    | "Deterministic"
    | "Statistical";

export type AITrainability = "None" | "Partial" | "Full";

export interface AIStats {
    speed: number;
    strategicDepth: number;
    adaptability: number;
    avgMoveTime?: string;
}

export interface AIDossier {
    id: string;
    name: string;
    tagline: string;
    description: string;
    image: string;
    concept: string;
    infoType: AIInfoType;
    trainability: AITrainability;
    trainabilityDetails: string;
    strengths: string[];
    weaknesses: string[];
    stats: AIStats;
    notes?: string;
}

export interface AITournamentBenchmark {
    aiId: string;
    aiName: string;
    sampleSize: string;
    totalRuntime: string;
    avgRoundsPerGame: string;
    flagCaptureRate: string;
    noMoveWinsRate: string;
    maxTurnCutoffsRate: string;
    notes: string;
}
