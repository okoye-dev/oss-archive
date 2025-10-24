import { apiService, getApiBaseUrl } from "./api-service";

export interface FileData {
  file_name: string;
  file_size: number;
  file_type: string;
  uploaded_at: string;
}

export const getFiles = async (): Promise<FileData[]> => {
  return apiService.get<FileData[]>("/files");
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
