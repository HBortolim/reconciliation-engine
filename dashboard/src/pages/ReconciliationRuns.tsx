import React, { useState, useEffect } from "react";
import { getRuns } from "../api/reconciliation";
import { ReconciliationRun } from "../types";
import StatusBadge from "../components/StatusBadge";
import MoneyDisplay from "../components/MoneyDisplay";

const ReconciliationRuns: React.FC = () => {
  const [runs, setRuns] = useState<ReconciliationRun[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchRuns = async () => {
      try {
        const data = await getRuns();
        setRuns(data);
      } catch (error) {
        console.error("Failed to fetch runs:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchRuns();
  }, []);

  if (loading) {
    return <div className="text-center text-gray-500">Loading...</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold text-gray-900">Reconciliation Runs</h1>
        <button className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
          New Run
        </button>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Run ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Start Time</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Processed</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Matched</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Exceptions</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Total Amount</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {runs.map((run) => (
                <tr key={run.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium text-blue-600 cursor-pointer hover:underline">
                    {run.id}
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <StatusBadge status={run.status} variant="status" />
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">
                    {new Date(run.startTime).toLocaleString()}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-600">{run.totalProcessed}</td>
                  <td className="px-6 py-4 text-sm text-green-600 font-medium">{run.matchedCount}</td>
                  <td className="px-6 py-4 text-sm text-red-600 font-medium">{run.exceptionCount}</td>
                  <td className="px-6 py-4 text-sm text-gray-900 font-mono">
                    <MoneyDisplay cents={run.statistics.totalAmount} />
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

export default ReconciliationRuns;
