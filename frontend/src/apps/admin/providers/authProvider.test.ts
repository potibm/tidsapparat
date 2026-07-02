import { describe, it, expect, vi, beforeEach } from "vitest";
import { configureOidc, authProvider, getAccessToken } from "./authProvider";

const mockUserManagerState = vi.hoisted(() => ({
  signinRedirect: vi.fn(),
  signinRedirectCallback: vi.fn(),
  getUser: vi.fn(),
  signoutRedirect: vi.fn(),
  removeUser: vi.fn(),
}));

vi.mock("oidc-client-ts", () => ({
  UserManager: vi.fn(function () {
    return mockUserManagerState;
  }),
}));

describe("authProvider", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    configureOidc("https://dex.example.com", "test-client");
  });

  describe("configureOidc", () => {
    it("creates a UserManager with correct settings", async () => {
      const { UserManager } = await import("oidc-client-ts");
      expect(UserManager).toHaveBeenCalledWith(
        expect.objectContaining({
          authority: "https://dex.example.com",
          client_id: "test-client",
          response_type: "code",
          scope: "openid profile email",
        }),
      );
    });
  });

  describe("login", () => {
    it("calls signinRedirect", async () => {
      await authProvider.login({});
      expect(mockUserManagerState.signinRedirect).toHaveBeenCalled();
    });
  });

  describe("handleCallback", () => {
    it("calls signinRedirectCallback", async () => {
      mockUserManagerState.signinRedirectCallback.mockResolvedValue({});
      await authProvider.handleCallback!();
      expect(mockUserManagerState.signinRedirectCallback).toHaveBeenCalled();
    });

    it("reuses active promise if called concurrently", async () => {
      mockUserManagerState.signinRedirectCallback.mockImplementation(
        () => new Promise((resolve) => setTimeout(resolve, 50)),
      );

      const p1 = authProvider.handleCallback!();
      const p2 = authProvider.handleCallback!();

      expect(await p1).toBe(await p2);
    });
  });

  describe("checkAuth", () => {
    it("resolves when user exists and is not expired", async () => {
      mockUserManagerState.getUser.mockResolvedValue({
        expired: false,
        access_token: "token123",
      });

      await expect(authProvider.checkAuth({})).resolves.toBeUndefined();
    });

    it("rejects when user is expired", async () => {
      mockUserManagerState.getUser.mockResolvedValue({ expired: true });

      await expect(authProvider.checkAuth({})).rejects.toThrow(
        "User not authenticated or token expired",
      );
    });

    it("rejects when no user", async () => {
      mockUserManagerState.getUser.mockResolvedValue(null);

      await expect(authProvider.checkAuth({})).rejects.toThrow(
        "User not authenticated or token expired",
      );
    });
  });

  describe("checkError", () => {
    it("rejects on 401", async () => {
      await expect(authProvider.checkError({ status: 401 })).rejects.toThrow(
        "Unauthorized API access",
      );
    });

    it("rejects on 403", async () => {
      await expect(authProvider.checkError({ status: 403 })).rejects.toThrow(
        "Unauthorized API access",
      );
    });

    it("resolves on other errors", async () => {
      await expect(
        authProvider.checkError({ status: 500 }),
      ).resolves.toBeUndefined();
    });
  });

  describe("logout", () => {
    it("calls signoutRedirect when user exists", async () => {
      mockUserManagerState.getUser.mockResolvedValue({ id: "user1" });
      mockUserManagerState.signoutRedirect.mockResolvedValue(undefined);

      await authProvider.logout({});
      expect(mockUserManagerState.signoutRedirect).toHaveBeenCalled();
    });

    it("calls removeUser when no user exists", async () => {
      mockUserManagerState.getUser.mockResolvedValue(null);

      await authProvider.logout({});
      expect(mockUserManagerState.removeUser).toHaveBeenCalled();
    });

    it("falls back to removeUser when end session endpoint is missing", async () => {
      mockUserManagerState.getUser.mockResolvedValue({ id: "user1" });
      mockUserManagerState.signoutRedirect.mockRejectedValue(
        new Error("No end session endpoint"),
      );

      await authProvider.logout({});
      expect(mockUserManagerState.removeUser).toHaveBeenCalled();
    });
  });

  describe("getIdentity", () => {
    it("returns identity from user profile", async () => {
      mockUserManagerState.getUser.mockResolvedValue({
        profile: {
          sub: "user-123",
          name: "Test User",
          picture: "https://example.com/avatar.png",
        },
      });

      const identity = await authProvider.getIdentity!();
      expect(identity.id).toBe("user-123");
      expect(identity.fullName).toBe("Test User");
      expect(identity.avatar).toBe("https://example.com/avatar.png");
    });

    it("uses preferred_username as fallback for fullName", async () => {
      mockUserManagerState.getUser.mockResolvedValue({
        profile: {
          sub: "user-123",
          preferred_username: "testuser",
        },
      });

      const identity = await authProvider.getIdentity!();
      expect(identity.fullName).toBe("testuser");
    });

    it("rejects when no profile", async () => {
      mockUserManagerState.getUser.mockResolvedValue({ profile: null });

      await expect(authProvider.getIdentity!()).rejects.toThrow(
        "No identity profile found",
      );
    });
  });

  describe("getPermissions", () => {
    it("returns roles from user profile", async () => {
      mockUserManagerState.getUser.mockResolvedValue({
        profile: { roles: ["admin", "editor"] },
      });

      const perms = await authProvider.getPermissions!({});
      expect(perms).toEqual(["admin", "editor"]);
    });

    it("returns empty array when no user", async () => {
      mockUserManagerState.getUser.mockResolvedValue(null);

      const perms = await authProvider.getPermissions!({});
      expect(perms).toEqual([]);
    });
  });

  describe("getAccessToken", () => {
    it("returns token when user is valid", async () => {
      mockUserManagerState.getUser.mockResolvedValue({
        expired: false,
        access_token: "valid-token",
      });

      const token = await getAccessToken();
      expect(token).toBe("valid-token");
    });

    it("returns null when user is expired", async () => {
      mockUserManagerState.getUser.mockResolvedValue({
        expired: true,
        access_token: "old",
      });

      const token = await getAccessToken();
      expect(token).toBeNull();
    });

    it("returns null when no user", async () => {
      mockUserManagerState.getUser.mockResolvedValue(null);

      const token = await getAccessToken();
      expect(token).toBeNull();
    });
  });
});
