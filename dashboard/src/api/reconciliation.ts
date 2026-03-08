import apiClient from "./client";
import {
  ReconciliationRun,
  ReconciliationException,
  AcquirerContract,
  ReconciliationPair
} from "../types";

export const getRuns = async (): Promise<ReconciliationRun[]> => {
  const response = await apiClient.get("/runs");
  return response.data;
};

export const getRunById = async (id: string): Promise<ReconciliationRun> => {
  const response = await apiClient.get(`/runs/${id}`);
  return response.data;
};

export const triggerRun = async (files: File[]): Promise<ReconciliationRun> => {
  const formData = new FormData();
  files.forEach((file) => {
    formData.append("files", file);
  });
  const response = await apiClient.post("/runs", formData, {
    headers: {
      "Content-Type": "multipart/form-data"
    }
  });
  return response.data;
};

export const getExceptions = async (
  runId: string
): Promise<ReconciliationException[]> => {
  const response = await apiClient.get(`/runs/${runId}/exceptions`);
  return response.data;
};

export const resolveException = async (
  id: string,
  note: string
): Promise<ReconciliationException> => {
  const response = await apiClient.patch(`/exceptions/${id}`, {
    resolutionStatus: "RESOLVED",
    resolutionNote: note
  });
  return response.data;
};

export const getContracts = async (): Promise<AcquirerContract[]> => {
  const response = await apiClient.get("/contracts");
  return response.data;
};

export const getFeeAnalysis = async (
  period: string
): Promise<Record<string, unknown>> => {
  const response = await apiClient.get(`/analysis/fees?period=${period}`);
  return response.data;
};

export const getPairs = async (runId: string): Promise<ReconciliationPair[]> => {
  const response = await apiClient.get(`/runs/${runId}/pairs`);
  return response.data;
};
