import { useState, useEffect, useCallback } from "react";
import { getFiles, uploadFiles, FileData } from "@/api/files";

interface UseQueryResult<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
  refetch: () => Promise<void>;
}

const useQuery = <T>(queryFn: () => Promise<T>): UseQueryResult<T> => {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await queryFn();
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
      console.error("Query error:", err);
    } finally {
      setLoading(false);
    }
  }, [queryFn]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  return {
    data,
    loading,
    error,
    refetch: fetchData,
  };
};

export const useFiles = () => {
  const { data: files, loading, error, refetch } = useQuery(() => getFiles());
  const [uploading, setUploading] = useState(false);

  const uploadMultipleFiles = useCallback(async (fileList: File[]): Promise<FileData[]> => {
    try {
      setUploading(true);
      const uploadedFiles: FileData[] = [];
      
      for (const file of fileList) {
        const formData = new FormData();
        formData.append("file", file);
        
        const uploadedFile = await uploadFiles(formData);
        uploadedFiles.push(uploadedFile);
      }

      await refetch(); // Refetch files after upload
      return uploadedFiles;
    } catch (err) {
      console.error("Failed to upload files:", err);
      return [];
    } finally {
      setUploading(false);
    }
  }, [refetch]);

  return {
    files: files || [],
    loading,
    error,
    uploading,
    refetch,
    uploadMultipleFiles,
  };
};