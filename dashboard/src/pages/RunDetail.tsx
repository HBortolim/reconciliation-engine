import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { getRunById, getPairs, getExceptions } from "../api/reconciliation";
import { ReconciliationRun, ReconciliationPair, ReconciliationException } from "../types";
import StatusBadge from "../components/StatusBadge";
import ConfidenceBar from "../components/ConfidenceBar";
import MoneyDisplay from "../components/MoneyDisplay";

const RunDetail: React.FC = () => {
  const { runId } = useParams<{ runId: string }>();
  const [run, setRun] = useState<ReconciliationRun | null>(null);
  const [pairs, setPairs] = useState<ReconciliationPair[]>([]);
  const [exceptions, setExceptions] = useState<ReconciliationException[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      if (!runId) return;

      try {
        const runData = await getRunById(runId);
        setRun(runData);

        const pairsData = await getPairs(runId);
        setPairs(pairsData);

        const exceptionsData = await getExceptions(runId);
        setExceptions(exceptionsData);
      } catch (error) {
        console.error("Failed to fetch run details:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [runId]);

  if (loading) {
    return <div className="text-center text-gray-500">Loading...</div>;
  }

  if (!run) {
    return <div className="text-center text-red-600">Run not found</div>;
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Run Details</h1>
        <p className="text-gray-500">{run.id}</p>
      </div>

      {/* Run Statistics */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium">Status</h3>
          <div className="mt-2">
            <StatusBadge status={run.status} variant="status" />
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium">Total Processed</h3>
          <p className="text-2xl font-bold text-gray-900 mt-2">{run.totalProcessed}</p>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium">Matched Rate</h3>
          <p className="text-2xl font-bold text-green-600 mt-2">
            {run.totalProcessed > 0 ? ((run.matchedCount / run.totalProcessed) * 100).toFixed(1) : 0}%
          </p>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-gray-500 text-sm font-medium">Total Amount</h3>
          <p className="text-lg font-mono font-bold text-gray-900 mt-2">
            <MoneyDisplay centavos={run.statistics.totalAmount} />
          </p>
        </div>
      </div>

      {/* Matched Pairs */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="p-6 border-b border-gray-200">
          <h2 className="text-xl font-bold text-gray-900">Matched Pairs ({pairs.length})</h2>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Pair ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Confidence</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Amount</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {pairs.slice(0, 10).map((pair) => (
                <tr key={pair.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">{pair.id}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{pair.matchType}</td>
                  <td className="px-6 py-4 text-sm">
                    <ConfidenceBar confidence={pair.confidence} />
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-900 font-mono">
                    <MoneyDisplay centavos={pair.transactionA.amount} />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Exceptions */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="p-6 border-b border-gray-200">
          <h2 className="text-xl font-bold text-gray-900">Exceptions ({exceptions.length})</h2>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Exception ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Severity</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {exceptions.slice(0, 10).map((exc) => (
                <tr key={exc.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">{exc.id}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{exc.exceptionType}</td>
                  <td className="px-6 py-4 text-sm">
                    <StatusBadge status={exc.severity} variant="severity" />
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <StatusBadge status={exc.resolutionStatus} variant="resolution" />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default RunDetail;
