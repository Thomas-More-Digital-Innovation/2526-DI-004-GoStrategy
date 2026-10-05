export type TrainingSupport =
    | "None"
    | "Partial"
    | "Full";

export interface AIDossier {
    id: string;
    name: string;
    category: string;
    description: string;
    image: string;
    algorithm: string;
    stateModel: string;
    complexity: string;
    trainingSupport: TrainingSupport;
    trainingDetails: string;
    strengths: string[];
    weaknesses: string[];
    notes?: string;
}

export interface AITournamentBenchmark {
    aiId: string;
    aiName: string;
    sampleSize: number;
    totalRuntimeSeconds: number;
    avgRoundsPerGame: number;
    flagCaptures: number;
    flagCaptureRate: string;
    noMoveWins: number;
    noMoveWinsRate: string;
    maxTurnCutoffs: number;
    maxTurnCutoffsRate: string;
    summary: string;
}
