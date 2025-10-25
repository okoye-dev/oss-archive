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
  return response.files || [];
};

export const uploadFiles = async (
  formData: FormData, 
  onProgress?: (progress: number) => void
): Promise<FileData> => {
  return new Promise((resolve, reject) => {
    const xhr = new XMLHttpRequest();
    
    if (onProgress) {
      xhr.upload.addEventListener('progress', (event) => {
        if (event.lengthComputable) {
          const progress = Math.round((event.loaded / event.total) * 100);
          onProgress(progress);
        }
      });
    }
    
    xhr.addEventListener('load', () => {
      if (xhr.status >= 200 && xhr.status < 300) {
        try {
          const result = JSON.parse(xhr.responseText);
          resolve(result);
        } catch (err) {
          reject(new Error('Failed to parse response'));
        }
      } else {
        reject(new Error(`Upload failed: ${xhr.statusText}`));
      }
    });
    
    xhr.addEventListener('error', () => {
      reject(new Error('Upload failed'));
    });
    
    xhr.open('POST', '/api/v1/files');
    xhr.send(formData);
  });
};

export const downloadFile = async (fileName: string): Promise<void> => {
  const response = await fetch(`/api/v1/files/${fileName}`);
  
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
