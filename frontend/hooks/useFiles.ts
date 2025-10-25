import { useState, useEffect, useCallback } from "react";
import { getFiles, uploadFiles, downloadFile, FileData } from "@/api/files";
import { useToast } from "./useToast";

export const useFiles = () => {
  const [files, setFiles] = useState<FileData[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const { toast } = useToast();

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
      setUploadProgress(0);
      const uploadedFiles: FileData[] = [];
      
      // Show initial toast
      const progressToast = toast({
        title: "Uploading files...",
        description: `0% complete (0/${fileList.length} files)`,
      });
      
      for (let i = 0; i < fileList.length; i++) {
        const file = fileList[i];
        const formData = new FormData();
        formData.append("file", file);
        
        const uploadedFile = await uploadFiles(formData, (progress) => {
          const fileProgress = (i / fileList.length) * 100 + (progress / fileList.length);
          const totalProgress = Math.round(fileProgress);
          setUploadProgress(totalProgress);
          
          // Update toast with progress
          progressToast.update({
            id: progressToast.id,
            title: "Uploading files...",
            description: `${totalProgress}% complete (${i + (progress === 100 ? 1 : 0)}/${fileList.length} files)`,
          });
        });
        uploadedFiles.push(uploadedFile);
      }

      // Show success toast
      progressToast.update({
        id: progressToast.id,
        title: "Upload complete! 🎉",
        description: `Successfully uploaded ${fileList.length} file${fileList.length > 1 ? 's' : ''}`,
      });

      await fetchFiles();
      return uploadedFiles;
    } catch (err) {
      console.error("Failed to upload files:", err);
      toast({
        title: "Upload failed",
        description: "Failed to upload files. Please try again.",
        variant: "destructive",
      });
      return [];
    } finally {
      setUploading(false);
      setUploadProgress(0);
    }
  }, [fetchFiles, toast]);

  const handleDownload = useCallback(async (fileName: string) => {
    try {
      await downloadFile(fileName);
      toast({
        title: "Download started",
        description: "File download has been initiated",
      });
    } catch (err) {
      console.error("Failed to download file:", err);
      toast({
        title: "Download failed",
        description: "Failed to download file. Please try again.",
        variant: "destructive",
      });
    }
  }, [toast]);

  return {
    files,
    loading,
    error,
    uploading,
    uploadProgress,
    refetchFiles: fetchFiles,
    uploadMultipleFiles,
    downloadFile: handleDownload,
  };
};