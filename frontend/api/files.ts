import { apiService, getApiBaseUrl } from "./api-service";

export interface FileData {
  id: string;
  name: string;
  storage_key: string;
}

interface FilesResponse {
  files: FileData[];
}

export const getFiles = async (): Promise<FileData[]> => {
  const response = await apiService.get<FilesResponse>("/files");
  return response.files;
};

export const uploadFiles = async (formData: FormData): Promise<FileData> => {
  const baseUrl = getApiBaseUrl();
  const response = await fetch(`${baseUrl}/api/v1/files`, {
    method: "POST",
    body: formData,
  });
  
  if (!response.ok) {
    throw new Error(`Failed to upload file: ${response.statusText}`);
  }
  
  return response.json();
};

export const downloadFile = async (fileName: string): Promise<void> => {
  const baseUrl = getApiBaseUrl();
  const response = await fetch(`${baseUrl}/v1/files/${fileName}`);
  
  if (!response.ok) {
    throw new Error(`Failed to download file: ${response.statusText}`);
  }
  
  const blob = await response.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = fileName;
  document.body.appendChild(a);
  a.click();
  window.URL.revokeObjectURL(url);
  document.body.removeChild(a);
};
