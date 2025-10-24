import { useState, useEffect, useCallback } from "react";
import { getFiles, uploadFiles, downloadFile, FileData } from "@/api/files";

export const useFiles = () => {
  const [files, setFiles] = useState<FileData[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [uploading, setUploading] = useState(false);

  const fetchFiles = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await getFiles();
      setFiles(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchFiles();
  }, []);

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

      await fetchFiles();
      return uploadedFiles;
    } catch (err) {
      console.error("Failed to upload files:", err);
      return [];
    } finally {
      setUploading(false);
    }
  }, [fetchFiles]);

  const handleDownload = useCallback(async (fileName: string) => {
    try {
      await downloadFile(fileName);
    } catch (err) {
      console.error("Failed to download file:", err);
    }
  }, []);

  return {
    files,
    loading,
    error,
    uploading,
    refetchFiles: fetchFiles,
    uploadMultipleFiles,
    downloadFile: handleDownload,
  };
};