import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { apiClient } from "../../core/api/client";
import MLModelDetailsPageView from "./MLModelDetailsPageView";

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

function MLModelDetailsPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const modelId = Number(id);

  const [name, setName] = useState("");
  const [modelData, setModelData] = useState("");
  const [fileName, setFileName] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isDragOver, setIsDragOver] = useState(false);
  const [showSuccessToast, setShowSuccessToast] = useState(false);

  const loadModel = async () => {
    if (!modelId) {
      setError("Некорректный ID модели.");
      return;
    }
    setIsLoading(true);
    setError(null);
    try {
      const response = await apiClient.api.getMlModel(modelId);
      setName(response.data.name);
      setModelData(response.data.model_data);
      setFileName(`${response.data.name}.onnx`);
    } catch (requestError) {
      setError("Не удалось загрузить ML-модель.");
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    void loadModel();
  }, [id]);

  const handleFilePicked = async (file: File) => {
    if (!file.name.toLowerCase().endsWith(".onnx")) {
      setError("Поддерживаются только файлы с расширением .onnx.");
      return;
    }
    try {
      const base64 = await readFileAsBase64(file);
      setModelData(base64);
      setFileName(file.name);
    } catch (requestError) {
      setError("Не удалось прочитать файл модели.");
    }
  };

  const handleSave = async () => {
    if (!modelId || !name.trim() || !modelData) {
      return;
    }
    setIsLoading(true);
    setError(null);
    try {
      await apiClient.api.updateMlModel(modelId, {
        name: name.trim(),
        model_data: modelData,
      });
      setShowSuccessToast(true);
    } catch (requestError) {
      setError("Не удалось сохранить ML-модель.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!modelId) {
      return;
    }
    setIsLoading(true);
    setError(null);
    try {
      await apiClient.api.deleteMlModel(modelId);
      navigate("/ml-models");
    } catch (requestError) {
      setError("Не удалось удалить ML-модель.");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <MLModelDetailsPageView
      title={name.trim() ? `ML-модель: ${name}` : "ML-модель"}
      name={name}
      modelSize={formatBytes(base64ByteLength(modelData))}
      fileName={fileName}
      isLoading={isLoading}
      error={error}
      onNameChange={setName}
      onSave={handleSave}
      onDelete={handleDelete}
      isDragOver={isDragOver}
      onDragOver={() => setIsDragOver(true)}
      onDragLeave={() => setIsDragOver(false)}
      onDrop={async (file) => {
        setIsDragOver(false);
        await handleFilePicked(file);
      }}
      onBrowseFile={handleFilePicked}
      showSuccessToast={showSuccessToast}
      onCloseSuccessToast={() => setShowSuccessToast(false)}
    />
  );
}

export default MLModelDetailsPage;
