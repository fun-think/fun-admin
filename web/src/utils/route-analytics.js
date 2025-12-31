import { reactive } from 'vue'

/**
 * 路由分析状态
 */
const analyticsState = reactive({
  // 页面访问统计
  pageViews: new Map(),

  // 路由切换记录
  routeHistory: [],

  // 性能数据
  performance: {
    loadTimes: new Map(),
    renderTimes: new Map(),
    errorRates: new Map(),
  },

  // 用户行为
  userBehavior: {
    bounceRate: 0,
    avgSessionTime: 0,
    pageDepth: new Map(),
  },

  // 错误统计
  errors: [],

  // 配置
  config: {
    enabled: true,
    maxHistorySize: 1000,
    maxErrorSize: 100,
    reportInterval: 30000, // 30秒上报一次
    enablePerformanceTracking: true,
    enableErrorTracking: true,
    enableUserBehaviorTracking: true,
  },
})

/**
 * 路由分析管理器
 */
export class RouteAnalytics {
  constructor(options = {}) {
    this.config = { ...analyticsState.config, ...options }
    this.sessionId = this.generateSessionId()
    this.startTime = Date.now()
    this.currentRoute = null
    this.routeStartTime = null

    // 定期上报数据
    if (this.config.reportInterval > 0) {
      this.reportTimer = setInterval(() => {
        this.reportAnalytics()
      }, this.config.reportInterval)
    }

    // 页面卸载时上报
    window.addEventListener('beforeunload', () => {
      this.reportAnalytics(true)
    })
  }

  /**
   * 生成会话ID
   */
  generateSessionId() {
    return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
  }

  /**
   * 记录路由进入
   * @param {object} to - 目标路由
   * @param {object} from - 来源路由
   */
  recordRouteEnter(to, from) {
    if (!this.config.enabled)
      return

    const now = Date.now()
    const routeKey = `${to.name || 'unknown'}:${to.path}`

    // 记录路由切换
    const routeChange = {
      sessionId: this.sessionId,
      timestamp: now,
      from: from
        ? {
            name: from.name,
            path: from.path,
            fullPath: from.fullPath,
          }
        : null,
      to: {
        name: to.name,
        path: to.path,
        fullPath: to.fullPath,
        meta: to.meta,
      },
      referrer: document.referrer,
      userAgent: navigator.userAgent,
    }

    // 添加到历史记录
    analyticsState.routeHistory.push(routeChange)

    // 限制历史记录大小
    if (analyticsState.routeHistory.length > this.config.maxHistorySize) {
      analyticsState.routeHistory.shift()
    }

    // 更新页面访问统计
    const current = analyticsState.pageViews.get(routeKey) || {
      count: 0,
      firstVisit: now,
      lastVisit: now,
      totalTime: 0,
      bounces: 0,
    }

    current.count++
    current.lastVisit = now
    analyticsState.pageViews.set(routeKey, current)

    // 记录当前路由和开始时间
    this.currentRoute = to
    this.routeStartTime = now

    // 性能追踪
    if (this.config.enablePerformanceTracking) {
      this.trackPerformance(to)
    }
  }

  /**
   * 记录路由离开
   * @param {object} from - 离开的路由
   * @param {object} to - 目标路由
   */
  recordRouteLeave(from, to) {
    if (!this.config.enabled || !this.routeStartTime)
      return

    const now = Date.now()
    const routeKey = `${from.name || 'unknown'}:${from.path}`
    const stayTime = now - this.routeStartTime

    // 更新停留时间
    const pageData = analyticsState.pageViews.get(routeKey)
    if (pageData) {
      pageData.totalTime += stayTime

      // 判断是否为跳出（停留时间少于30秒且没有进一步操作）
      if (stayTime < 30000 && this.isPageBounce(from)) {
        pageData.bounces++
      }

      analyticsState.pageViews.set(routeKey, pageData)
    }

    // 记录页面深度
    this.updatePageDepth(from, stayTime)
  }

  /**
   * 记录路由错误
   * @param {object} error - 错误信息
   * @param {object} route - 相关路由
   */
  recordRouteError(error, route) {
    if (!this.config.enabled || !this.config.enableErrorTracking)
      return

    const errorRecord = {
      sessionId: this.sessionId,
      timestamp: Date.now(),
      error: {
        message: error.message,
        stack: error.stack,
        name: error.name,
      },
      route: route
        ? {
            name: route.name,
            path: route.path,
            fullPath: route.fullPath,
          }
        : null,
      userAgent: navigator.userAgent,
      url: window.location.href,
    }

    analyticsState.errors.push(errorRecord)

    // 限制错误记录大小
    if (analyticsState.errors.length > this.config.maxErrorSize) {
      analyticsState.errors.shift()
    }

    // 更新错误率统计
    const routeKey = route ? `${route.name}:${route.path}` : 'unknown'
    const errorRate = analyticsState.performance.errorRates.get(routeKey) || { total: 0, errors: 0 }
    errorRate.errors++
    analyticsState.performance.errorRates.set(routeKey, errorRate)
  }

  /**
   * 性能追踪
   * @param {object} route - 路由对象
   */
  trackPerformance(route) {
    const routeKey = `${route.name}:${route.path}`

    // 记录加载开始时间
    const loadStartTime = Date.now()

    // 监听页面加载完成
    if (document.readyState === 'loading') {
      document.addEventListener('DOMContentLoaded', () => {
        this.recordLoadTime(routeKey, Date.now() - loadStartTime)
      })
    }
    else {
      // 页面已加载，记录渲染时间
      this.$nextTick(() => {
        this.recordRenderTime(routeKey, Date.now() - loadStartTime)
      })
    }

    // 使用Performance API获取更精确的数据
    if (window.performance && window.performance.getEntriesByType) {
      setTimeout(() => {
        const navigationEntries = window.performance.getEntriesByType('navigation')
        if (navigationEntries.length > 0) {
          const entry = navigationEntries[0]
          this.recordDetailedPerformance(routeKey, entry)
        }
      }, 100)
    }
  }

  /**
   * 记录加载时间
   * @param {string} routeKey - 路由键
   * @param {number} loadTime - 加载时间
   */
  recordLoadTime(routeKey, loadTime) {
    const times = analyticsState.performance.loadTimes.get(routeKey) || []
    times.push(loadTime)
    analyticsState.performance.loadTimes.set(routeKey, times)
  }

  /**
   * 记录渲染时间
   * @param {string} routeKey - 路由键
   * @param {number} renderTime - 渲染时间
   */
  recordRenderTime(routeKey, renderTime) {
    const times = analyticsState.performance.renderTimes.get(routeKey) || []
    times.push(renderTime)
    analyticsState.performance.renderTimes.set(routeKey, times)
  }

  /**
   * 记录详细性能数据
   * @param {string} routeKey - 路由键
   * @param {object} entry - 性能条目
   */
  recordDetailedPerformance(routeKey, entry) {
    const performanceData = {
      dns: entry.domainLookupEnd - entry.domainLookupStart,
      tcp: entry.connectEnd - entry.connectStart,
      ssl: entry.secureConnectionStart > 0 ? entry.connectEnd - entry.secureConnectionStart : 0,
      ttfb: entry.responseStart - entry.requestStart,
      download: entry.responseEnd - entry.responseStart,
      dom: entry.domContentLoadedEventEnd - entry.domContentLoadedEventStart,
      load: entry.loadEventEnd - entry.loadEventStart,
    }

    // 存储详细性能数据
    const existing = analyticsState.performance.loadTimes.get(`${routeKey}:detailed`) || []
    existing.push(performanceData)
    analyticsState.performance.loadTimes.set(`${routeKey}:detailed`, existing)
  }

  /**
   * 判断是否为页面跳出
   * @param {object} route - 路由对象
   */
  isPageBounce(route) {
    // 简单的跳出判断：用户没有进行任何交互就离开
    // 这里可以根据具体需求调整判断逻辑
    return analyticsState.routeHistory.length <= 1
  }

  /**
   * 更新页面深度
   * @param {object} route - 路由对象
   * @param {number} stayTime - 停留时间
   */
  updatePageDepth(route, stayTime) {
    const routeKey = `${route.name}:${route.path}`
    const current = analyticsState.userBehavior.pageDepth.get(routeKey) || {
      totalViews: 0,
      totalTime: 0,
      avgTime: 0,
    }

    current.totalViews++
    current.totalTime += stayTime
    current.avgTime = current.totalTime / current.totalViews

    analyticsState.userBehavior.pageDepth.set(routeKey, current)
  }

  /**
   * 获取分析报告
   */
  getAnalyticsReport() {
    return {
      sessionId: this.sessionId,
      sessionDuration: Date.now() - this.startTime,
      pageViews: this.getPageViewsReport(),
      performance: this.getPerformanceReport(),
      userBehavior: this.getUserBehaviorReport(),
      errors: this.getErrorsReport(),
      routeHistory: analyticsState.routeHistory.slice(-50), // 最近50条记录
    }
  }

  /**
   * 获取页面访问报告
   */
  getPageViewsReport() {
    const report = []
    analyticsState.pageViews.forEach((data, routeKey) => {
      report.push({
        route: routeKey,
        ...data,
        avgStayTime: data.totalTime / data.count,
        bounceRate: data.bounces / data.count,
      })
    })
    return report.sort((a, b) => b.count - a.count)
  }

  /**
   * 获取性能报告
   */
  getPerformanceReport() {
    const report = {
      loadTimes: {},
      renderTimes: {},
      errorRates: {},
    }

    // 加载时间统计
    analyticsState.performance.loadTimes.forEach((times, routeKey) => {
      if (Array.isArray(times)) {
        const avg = times.reduce((a, b) => a + b, 0) / times.length
        const min = Math.min(...times)
        const max = Math.max(...times)
        report.loadTimes[routeKey] = { avg, min, max, count: times.length }
      }
    })

    // 渲染时间统计
    analyticsState.performance.renderTimes.forEach((times, routeKey) => {
      if (Array.isArray(times)) {
        const avg = times.reduce((a, b) => a + b, 0) / times.length
        const min = Math.min(...times)
        const max = Math.max(...times)
        report.renderTimes[routeKey] = { avg, min, max, count: times.length }
      }
    })

    // 错误率统计
    analyticsState.performance.errorRates.forEach((data, routeKey) => {
      const pageData = analyticsState.pageViews.get(routeKey)
      if (pageData) {
        data.total = pageData.count
        data.errorRate = data.errors / data.total
      }
      report.errorRates[routeKey] = data
    })

    return report
  }

  /**
   * 获取用户行为报告
   */
  getUserBehaviorReport() {
    const totalPageViews = Array.from(analyticsState.pageViews.values())
      .reduce((sum, data) => sum + data.count, 0)

    const totalBounces = Array.from(analyticsState.pageViews.values())
      .reduce((sum, data) => sum + data.bounces, 0)

    const totalTime = Array.from(analyticsState.pageViews.values())
      .reduce((sum, data) => sum + data.totalTime, 0)

    return {
      bounceRate: totalBounces / totalPageViews,
      avgSessionTime: totalTime / totalPageViews,
      totalPageViews,
      uniqueRoutes: analyticsState.pageViews.size,
      pageDepth: Object.fromEntries(analyticsState.userBehavior.pageDepth),
    }
  }

  /**
   * 获取错误报告
   */
  getErrorsReport() {
    const errorsByRoute = {}
    const errorsByType = {}

    analyticsState.errors.forEach((error) => {
      const routeKey = error.route ? `${error.route.name}:${error.route.path}` : 'unknown'
      const errorType = error.error.name || 'Unknown'

      errorsByRoute[routeKey] = (errorsByRoute[routeKey] || 0) + 1
      errorsByType[errorType] = (errorsByType[errorType] || 0) + 1
    })

    return {
      total: analyticsState.errors.length,
      errorsByRoute,
      errorsByType,
      recentErrors: analyticsState.errors.slice(-10), // 最近10个错误
    }
  }

  /**
   * 上报分析数据
   * @param {boolean} immediate - 是否立即上报
   */
  async reportAnalytics(immediate = false) {
    if (!this.config.enabled)
      return

    const report = this.getAnalyticsReport()

    try {
      // 发送到服务器
      if (this.config.reportEndpoint) {
        await fetch(this.config.reportEndpoint, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
          },
          body: JSON.stringify(report),
        })
      }

      // 本地存储
      if (this.config.enableLocalStorage) {
        localStorage.setItem('route_analytics', JSON.stringify(report))
      }

      console.log('Analytics reported successfully:', report)
    }
    catch (error) {
      console.error('Failed to report analytics:', error)
    }
  }

  /**
   * 清理数据
   */
  cleanup() {
    if (this.reportTimer) {
      clearInterval(this.reportTimer)
    }
  }

  /**
   * 获取状态
   */
  getState() {
    return analyticsState
  }
}

/**
 * 路由分析组合式API
 */
export function useRouteAnalytics(options = {}) {
  const analytics = new RouteAnalytics(options)

  return {
    // 状态
    state: analyticsState,

    // 方法
    recordRouteEnter: (to, from) => analytics.recordRouteEnter(to, from),
    recordRouteLeave: (from, to) => analytics.recordRouteLeave(from, to),
    recordRouteError: (error, route) => analytics.recordRouteError(error, route),
    getReport: () => analytics.getAnalyticsReport(),
    reportAnalytics: () => analytics.reportAnalytics(),
    cleanup: () => analytics.cleanup(),

    // 实例
    analytics,
  }
}

// 全局分析器实例
let globalAnalytics = null

/**
 * 创建全局路由分析器
 * @param {object} options - 配置选项
 */
export function createRouteAnalytics(options = {}) {
  globalAnalytics = new RouteAnalytics(options)
  return globalAnalytics
}

/**
 * 获取全局路由分析器
 */
export function getRouteAnalytics() {
  return globalAnalytics
}

export default RouteAnalytics
