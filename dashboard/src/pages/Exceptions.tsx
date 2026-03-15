import React, { useState, useEffect } from "react";
import { getExceptions, resolveException } from "../api/reconciliation";
import { ReconciliationException, ExceptionType, Severity } from "../types";
import StatusBadge from "../components/StatusBadge";
import MoneyDisplay from "../components/MoneyDisplay";

const Exceptions: React.FC = () => {
  const [exceptions, setExceptions] = useState<ReconciliationException[]>([]);
  const [loading, setLoading] = useState(true);
  const [filterSeverity, setFilterSeverity] = useState<string | null>(null);
  const [selectedExceptionId, setSelectedExceptionId] = useState<string | null>(null);
  const [resolutionNote, setResolutionNote] = useState("");

  useEffect(() => {
    const fetchExceptions = async () => {
      try {
        // Fetch from first run as example
        const response = await getExceptions("default-run-id");
        setExceptions(response);
      } catch (error) {
        console.error("Failed to fetch exceptions:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchExceptions();
  }, []);

  const handleResolve = async () => {
    if (!selectedExceptionId) return;

    try {
      await resolveException(selectedExceptionId, resolutionNote);
      setExceptions(exceptions.map(exc =>
        exc.id === selectedExceptionId
          ? { ...exc, resolutionStatus: "RESOLVED" as const, resolutionNote }
          : exc
      ));
      setSelectedExceptionId(null);
      setResolutionNote("");
    } catch (error) {
      console.error("Failed to resolve exception:", error);
    }
  };

  const filteredExceptions = filterSeverity
    ? exceptions.filter(exc => exc.severity === filterSeverity)
    : exceptions;

  if (loading) {
    return <div className="text-center text-gray-500">Loading...</div>;
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-900">Exception Management</h1>

      {/* Filters */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="space-y-4">
          <label className="block">
            <span className="text-sm font-medium text-gray-700">Filter by Severity</span>
            <select
              value={filterSeverity || ""}
              onChange={(e) => setFilterSeverity(e.target.value || null)}
              className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-lg text-gray-900"
            >
              <option value="">All Severities</option>
              <option value={Severity.LOW}>Low</option>
              <option value={Severity.MEDIUM}>Medium</option>
              <option value={Severity.HIGH}>High</option>
              <option value={Severity.CRITICAL}>Critical</option>
            </select>
          </label>
        </div>
      </div>

      {/* Exceptions Table */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-gray-50 border-b border-gray-200">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Exception ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Severity</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Discrepancy</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Created</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {filteredExceptions.map((exc) => (
                <tr key={exc.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-medium text-gray-900">{exc.id}</td>
                  <td className="px-6 py-4 text-sm text-gray-600">{exc.exceptionType}</td>
                  <td className="px-6 py-4 text-sm">
                    <StatusBadge status={exc.severity} variant="severity" />
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <StatusBadge status={exc.resolutionStatus} variant="resolution" />
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-900 font-mono">
                    {exc.discrepancyAmount ? (
                      <MoneyDisplay cents={exc.discrepancyAmount} />
                    ) : (
                      "-"
                    )}
                  </td>
                  <td className="px-6 py-4 text-sm text-gray-500">
                    {new Date(exc.createdAt).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <button
                      onClick={() => setSelectedExceptionId(exc.id)}
                      className="text-blue-600 hover:text-blue-800 font-medium"
                    >
                      Review
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Resolution Modal */}
      {selectedExceptionId && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white rounded-lg shadow-lg max-w-md w-full p-6">
            <h3 className="text-lg font-bold text-gray-900 mb-4">Resolve Exception</h3>
            <textarea
              value={resolutionNote}
              onChange={(e) => setResolutionNote(e.target.value)}
              placeholder="Add resolution notes..."
              className="w-full px-3 py-2 border border-gray-300 rounded-lg text-gray-900 mb-4"
              rows={4}
            />
            <div className="flex gap-3">
              <button
                onClick={handleResolve}
                className="flex-1 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
              >
                Resolve
              </button>
              <button
                onClick={() => setSelectedExceptionId(null)}
                className="flex-1 px-4 py-2 bg-gray-300 text-gray-900 rounded-lg hover:bg-gray-400 transition-colors"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Exceptions;
