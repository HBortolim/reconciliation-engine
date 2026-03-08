// Enums
export enum SourceType {
  ACQUIRER = "ACQUIRER",
  ISSUER = "ISSUER",
  GATEWAY = "GATEWAY"
}

export enum MatchType {
  EXACT = "EXACT",
  FUZZY = "FUZZY",
  MANUAL = "MANUAL"
}

export enum ExceptionType {
  AMOUNT_MISMATCH = "AMOUNT_MISMATCH",
  MISSING_TRANSACTION = "MISSING_TRANSACTION",
  DUPLICATE = "DUPLICATE",
  DATE_MISMATCH = "DATE_MISMATCH",
  FEE_DISCREPANCY = "FEE_DISCREPANCY"
}

export enum Severity {
  LOW = "LOW",
  MEDIUM = "MEDIUM",
  HIGH = "HIGH",
  CRITICAL = "CRITICAL"
}

export enum RunStatus {
  PENDING = "PENDING",
  RUNNING = "RUNNING",
  COMPLETED = "COMPLETED",
  FAILED = "FAILED"
}

export enum ResolutionStatus {
  UNRESOLVED = "UNRESOLVED",
  PENDING_REVIEW = "PENDING_REVIEW",
  RESOLVED = "RESOLVED",
  ESCALATED = "ESCALATED"
}

// Domain Types
export interface TransactionRecord {
  id: string;
  sourceType: SourceType;
  transactionId: string;
  amount: number; // in centavos
  date: string;
  description: string;
  acquirerId: string;
  bandeira: string;
  authorizationCode?: string;
  metadata: Record<string, unknown>;
}

export interface ReconciliationPair {
  id: string;
  runId: string;
  transactionA: TransactionRecord;
  transactionB: TransactionRecord;
  matchType: MatchType;
  confidence: number; // 0-100
  matchedAt: string;
}

export interface ReconciliationException {
  id: string;
  runId: string;
  exceptionType: ExceptionType;
  severity: Severity;
  resolutionStatus: ResolutionStatus;
  primaryTransaction: TransactionRecord;
  relatedTransactions?: TransactionRecord[];
  discrepancyAmount?: number; // in centavos
  createdAt: string;
  resolvedAt?: string;
  resolutionNote?: string;
  assignedTo?: string;
}

export interface ReconciliationRun {
  id: string;
  status: RunStatus;
  startTime: string;
  endTime?: string;
  totalProcessed: number;
  matchedCount: number;
  exceptionCount: number;
  pendingReviewCount: number;
  sourceFiles: string[];
  statistics: {
    totalAmount: number;
    matchedAmount: number;
    discrepancyAmount: number;
  };
}

export interface AcquirerContract {
  id: string;
  acquirerId: string;
  bandeira: string;
  effectiveDate: string;
  expiryDate?: string;
  feeStructure: FeeSchedule[];
  status: "ACTIVE" | "INACTIVE" | "PENDING";
}

export interface FeeSchedule {
  id: string;
  acquirerContractId: string;
  feeType: string; // e.g., "INTERCHANGE", "MDR", "PROCESSING"
  feeValue: number; // percentage or fixed amount
  feeBase: "PERCENTAGE" | "FIXED";
  applicableRange?: {
    minAmount?: number;
    maxAmount?: number;
  };
}
