import { useMutation } from "@/api/helpers";
import { authService } from "@/api/auth-service";
import { SignupData, useAuthStore } from "@/store/authStore";

/**
 * Custom hook for farmer signup with validation and error handling
 */
export const useSignup = () => {
  const { signUp } = useAuthStore();

  return useMutation(
    async (data: SignupData) => {
      const result = await authService.signup(data);
      return result;
    },
    {
      maxRetries: 2,
      retryDelay: 1000,
    }
  );
};

/**
 * Custom hook for farmer signin with validation and error handling
 */
export const useSignin = () => {
  const { signIn } = useAuthStore();

  return useMutation(
    async (credentials: { username: string; password: string }) => {
      const result = await authService.signin(
        credentials.username,
        credentials.password
      );
      return result;
    },
    {
      maxRetries: 2,
      retryDelay: 1000,
    }
  );
};

/**
 * Custom hook to get current auth state
 */
export const useAuthState = () => {
  const { user, isAuthenticated, isLoading, error, clearError, signOut } =
    useAuthStore();

  return {
    user,
    isAuthenticated,
    isLoading,
    error,
    clearError,
    signOut,
  };
};

