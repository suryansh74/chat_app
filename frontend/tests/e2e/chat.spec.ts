import { test, expect, Page } from '@playwright/test';

const API_BASE = 'http://localhost:8000/api';

function testUser(index: number) {
  return {
    name: `TestUser${index}`,
    email: `testuser${index}@example.com`,
    password: 'TestPassword123!',
  };
}

async function registerUser(page: Page, index: number) {
  const user = testUser(index);
  await page.goto('/register');
  await page.getByLabel(/name/i).fill(user.name);
  await page.getByLabel(/email/i).fill(user.email);
  await page.getByLabel(/password/i).fill(user.password);
  await page.getByLabel(/confirm password/i).fill(user.password);
  await page.getByRole('button', { name: /register/i }).click();
  await page.waitForURL('/verify-email', { timeout: 10000 });
  return user;
}

async function loginUser(page: Page, index: number) {
  const user = testUser(index);
  await page.goto('/login');
  await page.getByLabel(/email/i).fill(user.email);
  await page.getByLabel(/password/i).fill(user.password);
  await page.getByRole('button', { name: /login/i }).click();
  await page.waitForURL('/verify-email', { timeout: 10000 });
  return user;
}

async function verifyEmail(page: Page, index: number) {
  await page.waitForTimeout(1000);

  await page.goto('/verify-email');
  await page.getByRole('button', { name: /send/i }).click();
  await page.waitForTimeout(2000);

  const response = await page.request.get(`${API_BASE}/email_verification/verified`);
  await expect(response.ok()).toBeTruthy();

  const profile = await page.request.get(`${API_BASE}/profile`);
  const profileData = await profile.json();
  const userId = profileData.data.user.id;

  return userId;
}

async function getUserId(page: Page): Promise<string> {
  const response = await page.request.get(`${API_BASE}/profile`);
  const data = await response.json();
  return data.data.user.id;
}

async function sendFriendRequest(page: Page, email: string) {
  await page.getByTitle('Add Friend').click();
  await page.getByPlaceholder(/email/i).fill(email);
  await page.getByRole('button', { name: /send/i }).click();
  await page.waitForTimeout(1000);
}

test.describe('Real-time Chat Features', () => {
  test('Friend request sends real-time notification badge to recipient', async ({ browser }) => {
    const context1 = await browser.newContext();
    const context2 = await browser.newContext();
    const page1 = await context1.newPage();
    const page2 = await context2.newPage();

    await registerUser(page1, 100);
    await registerUser(page2, 101);

    await page1.goto('/verify-email');
    await page2.goto('/verify-email');

    await page1.getByRole('button', { name: /send/i }).click();
    await page2.getByRole('button', { name: /send/i }).click();
    await page1.waitForTimeout(2000);
    await page2.waitForTimeout(2000);

    await page1.goto('/home');
    await page2.goto('/home');
    await page1.waitForTimeout(2000);
    await page2.waitForTimeout(2000);

    const user1Id = await getUserId(page1);
    const user2Id = await getUserId(page2);

    await sendFriendRequest(page1, testUser(101).email);

    const badge = page2.locator('button[title="Notifications"] span.absolute');
    await expect(badge).toBeVisible({ timeout: 10000 });
    await expect(badge).toContainText('1');

    await context1.close();
    await context2.close();
  });

  test('Accept friend request updates friend list in real-time', async ({ browser }) => {
    const context1 = await browser.newContext();
    const context2 = await browser.newContext();
    const page1 = await context1.newPage();
    const page2 = await context2.newPage();

    await registerUser(page1, 200);
    await registerUser(page2, 201);

    await page1.goto('/verify-email');
    await page2.goto('/verify-email');
    await page1.getByRole('button', { name: /send/i }).click();
    await page2.getByRole('button', { name: /send/i }).click();
    await page1.waitForTimeout(2000);
    await page2.waitForTimeout(2000);

    await page1.goto('/home');
    await page2.goto('/home');
    await page1.waitForTimeout(2000);
    await page2.waitForTimeout(2000);

    await sendFriendRequest(page1, testUser(201).email);

    await page2.locator('button[title="Notifications"]').click();
    await page2.waitForTimeout(1000);
    await page2.getByRole('button', { name: /accept/i }).click();
    await page2.waitForTimeout(2000);

    await page1.getByRole('button', { name: /notifications/i }).click();
    await page1.getByRole('button', { name: /accept/i }).click();
    await page1.waitForTimeout(2000);

    const friendItem1 = page1.locator('text=' + testUser(201).name).first();
    await expect(friendItem1).toBeVisible({ timeout: 10000 });

    const friendItem2 = page2.locator('text=' + testUser(200).name).first();
    await expect(friendItem2).toBeVisible({ timeout: 10000 });

    await context1.close();
    await context2.close();
  });

  test('Messages are received in real-time by recipient', async ({ browser }) => {
    const context1 = await browser.newContext();
    const context2 = await browser.newContext();
    const page1 = await context1.newPage();
    const page2 = await context2.newPage();

    await registerUser(page1, 300);
    await registerUser(page2, 301);

    await page1.goto('/verify-email');
    await page2.goto('/verify-email');
    await page1.getByRole('button', { name: /send/i }).click();
    await page2.getByRole('button', { name: /send/i }).click();
    await page1.waitForTimeout(2000);
    await page2.waitForTimeout(2000);

    await page1.goto('/home');
    await page2.goto('/home');
    await page1.waitForTimeout(2000);
    await page2.waitForTimeout(2000);

    await sendFriendRequest(page1, testUser(301).email);

    await page2.locator('button[title="Notifications"]').click();
    await page2.waitForTimeout(1000);
    await page2.getByRole('button', { name: /accept/i }).click();
    await page2.waitForTimeout(2000);

    await page1.locator('button[title="Notifications"]').click();
    await page1.waitForTimeout(1000);
    await page1.getByRole('button', { name: /accept/i }).click();
    await page1.waitForTimeout(2000);

    await page1.getByText(testUser(301).name).click();
    await page2.getByText(testUser(300).name).click();
    await page1.waitForTimeout(1000);
    await page2.waitForTimeout(1000);

    await page1.locator('input[placeholder="Type a message..."]').fill('Hello from user 1!');
    await page1.getByRole('button', { name: /send/i }).click();

    const receivedMessage = page2.locator('text=Hello from user 1!');
    await expect(receivedMessage).toBeVisible({ timeout: 10000 });

    await page2.locator('input[placeholder="Type a message..."]').fill('Hello back!');
    await page2.getByRole('button', { name: /send/i }).click();

    const replyMessage = page1.locator('text=Hello back!');
    await expect(replyMessage).toBeVisible({ timeout: 10000 });

    await context1.close();
    await context2.close();
  });

  test('Unread count does not increment when chat is open', async ({ browser }) => {
    const context1 = await browser.newContext();
    const context2 = await browser.newContext();
    const page1 = await context1.newPage();
    const page2 = await context2.newPage();

    await registerUser(page1, 400);
    await registerUser(page2, 401);

    await page1.goto('/verify-email');
    await page2.goto('/verify-email');
    await page1.getByRole('button', { name: /send/i }).click();
    await page2.getByRole('button', { name: /send/i }).click();
    await page1.waitForTimeout(2000);
    await page2.waitForTimeout(2000);

    await page1.goto('/home');
    await page2.goto('/home');
    await page1.waitForTimeout(2000);
    await page2.waitForTimeout(2000);

    await sendFriendRequest(page1, testUser(401).email);

    await page2.locator('button[title="Notifications"]').click();
    await page2.waitForTimeout(1000);
    await page2.getByRole('button', { name: /accept/i }).click();
    await page2.waitForTimeout(2000);

    await page1.locator('button[title="Notifications"]').click();
    await page1.waitForTimeout(1000);
    await page1.getByRole('button', { name: /accept/i }).click();
    await page1.waitForTimeout(2000);

    await page1.getByText(testUser(401).name).click();
    await page2.getByText(testUser(400).name).click();
    await page1.waitForTimeout(1000);
    await page2.waitForTimeout(1000);

    await page1.locator('input[placeholder="Type a message..."]').fill('Test unread count!');
    await page1.getByRole('button', { name: /send/i }).click();
    await page2.waitForTimeout(2000);

    const unreadBadge = page2.locator('button[title="Notifications"] span.absolute');
    await expect(unreadBadge).not.toBeVisible();

    await context1.close();
    await context2.close();
  });
});
