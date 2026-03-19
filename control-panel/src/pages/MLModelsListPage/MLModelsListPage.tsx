import { useEffect, useMemo, useState } from "react";
import { apiClient } from "../../core/api/client";
import MLModelsListPageView, { type MLModelRow } from "./MLModelsListPageView";

function readFileAsBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      if (typeof reader.result !== "string") {
        reject(new Error("invalid file result"));
        return;
      }
      const payload = reader.result.includes(",")
        ? reader.result.split(",")[1]
        : reader.result;
      resolve(payload);
    };
    reader.onerror = () => reject(new Error("failed to read file"));
    reader.readAsDataURL(file);
  });
}

function base64ByteLength(value: string): number {
  const normalized = value.replace(/=+$/, "");
  return Math.floor((normalized.length * 3) / 4);
}

function formatBytes(bytes: number): string {
  if (bytes < 1024) {
    return `${bytes} B`;
  }
  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
}

function sanitizeModelName(fileName: string): string {
  return fileName.replace(/\.onnx$/i, "").trim();
}

function MLModelsListPage() {
  const [models, setModels] = useState<MLModelRow[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [nameFilter, setNameFilter] = useState("");

  const [showCreateModal, setShowCreateModal] = useState(false);
  const [createName, setCreateName] = useState("");
  const [createModelData, setCreateModelData] = useState("");
  const [createFileName, setCreateFileName] = useState("");
  const [isDragOver, setIsDragOver] = useState(false);
  const [createAttempted, setCreateAttempted] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const loadModels = async () => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.api.listMlModels();
      setModels(
        response.data.map((model) => ({
          id: model.id,
          name: model.name,
          sizeBytes: base64ByteLength(model.model_data),
        })),
      );
    } catch (requestError) {
      setError("Не удалось загрузить ML-модели.");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    void loadModels();
  }, []);

  const filteredModels = useMemo(
    () =>
      models
        .filter((model) =>
          model.name.toLowerCase().includes(nameFilter.toLowerCase()),
        )
        .sort((a, b) => a.name.localeCompare(b.name)),
    [models, nameFilter],
  );

  useEffect(() => {
    setCurrentPage(1);
  }, [nameFilter]);

  const totalPages = Math.max(1, Math.ceil(filteredModels.length / pageSize));
  const normalizedCurrentPage = Math.min(currentPage, totalPages);

  const paginatedModels = useMemo(() => {
    const start = (normalizedCurrentPage - 1) * pageSize;
    return filteredModels.slice(start, start + pageSize);
  }, [filteredModels, normalizedCurrentPage, pageSize]);

  const handleFilePicked = async (file: File) => {
    if (!file.name.toLowerCase().endsWith(".onnx")) {
      setError("Поддерживаются только файлы с расширением .onnx.");
      return;
    }
    try {
      const base64 = await readFileAsBase64(file);
      setCreateModelData(base64);
      setCreateFileName(file.name);
      if (!createName.trim()) {
        setCreateName(sanitizeModelName(file.name));
      }
    } catch (requestError) {
      setError("Не удалось прочитать файл модели.");
    }
  };

  const handleCreateSubmit = async () => {
    setCreateAttempted(true);
    if (!createName.trim() || !createModelData) {
      return;
    }

    setError(null);
    setIsLoading(true);
    try {
      await apiClient.api.createMlModel({
        name: createName.trim(),
        model_data: createModelData,
      });
      setShowCreateModal(false);
      setCreateAttempted(false);
      setCreateName("");
      setCreateModelData("");
      setCreateFileName("");
      await loadModels();
    } catch (requestError) {
      setError("Не удалось создать ML-модель.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <MLModelsListPageView
      models={paginatedModels}
      formatBytes={formatBytes}
      isLoading={isLoading}
      error={error}
      nameFilter={nameFilter}
      onNameFilterChange={setNameFilter}
      showCreateModal={showCreateModal}
      onOpenCreateModal={() => setShowCreateModal(true)}
      onCloseCreateModal={() => {
        setShowCreateModal(false);
        setCreateAttempted(false);
      }}
      createName={createName}
      onCreateNameChange={setCreateName}
      createFileName={createFileName}
      createModelData={createModelData}
      onCreateSubmit={handleCreateSubmit}
      createAttempted={createAttempted}
      isDragOver={isDragOver}
      onDragOver={() => setIsDragOver(true)}
      onDragLeave={() => setIsDragOver(false)}
      onDrop={async (file) => {
        setIsDragOver(false);
        await handleFilePicked(file);
      }}
      onBrowseFile={handleFilePicked}
      currentPage={normalizedCurrentPage}
      totalPages={totalPages}
      onPageChange={setCurrentPage}
      pageSize={pageSize}
      onPageSizeChange={(value) => {
        setPageSize(value);
        setCurrentPage(1);
      }}
    />
  );
}

export default MLModelsListPage;
