import React from "react";
import { RunStatus, ResolutionStatus, ExceptionType, Severity } from "../types";

interface StatusBadgeProps {
  status:
    | RunStatus
    | ResolutionStatus
    | ExceptionType
    | Severity
    | string;
  variant?: "status" | "resolution" | "exception" | "severity";
}

const getColorClasses = (status: string, variant?: string): string => {
  if (variant === "status" || variant === "resolution") {
    switch (status) {
      case RunStatus.PENDING:
      case ResolutionStatus.PENDING_REVIEW:
        return "bg-yellow-100 text-yellow-800 border border-yellow-300";
      case RunStatus.RUNNING:
        return "bg-blue-100 text-blue-800 border border-blue-300";
      case RunStatus.COMPLETED:
      case ResolutionStatus.RESOLVED:
        return "bg-green-100 text-green-800 border border-green-300";
      case RunStatus.FAILED:
        return "bg-red-100 text-red-800 border border-red-300";
      case ResolutionStatus.ESCALATED:
        return "bg-orange-100 text-orange-800 border border-orange-300";
      case ResolutionStatus.UNRESOLVED:
        return "bg-gray-100 text-gray-800 border border-gray-300";
      default:
        return "bg-gray-100 text-gray-800 border border-gray-300";
    }
  }

  if (variant === "severity") {
    switch (status) {
      case Severity.LOW:
        return "bg-green-100 text-green-800 border border-green-300";
      case Severity.MEDIUM:
        return "bg-yellow-100 text-yellow-800 border border-yellow-300";
      case Severity.HIGH:
        return "bg-orange-100 text-orange-800 border border-orange-300";
      case Severity.CRITICAL:
        return "bg-red-100 text-red-800 border border-red-300";
      default:
        return "bg-gray-100 text-gray-800 border border-gray-300";
    }
  }

  return "bg-gray-100 text-gray-800 border border-gray-300";
};

const StatusBadge: React.FC<StatusBadgeProps> = ({ status, variant }) => {
  return (
    <span
      className={`px-3 py-1 rounded-full text-sm font-medium ${getColorClasses(
        status,
        variant
      )}`}
    >
      {String(status).replace(/_/g, " ")}
    </span>
  );
};

export default StatusBadge;
