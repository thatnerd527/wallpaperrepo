using Microsoft.JSInterop;

using WallpaperUI.Pages;

namespace WallpaperUI.Cs.API
{
    public class PanelManagement
    {
        [JSInvokable]
        public static void ClosePanel(string fixedpanelid)
        {
            Home.panels.RemoveAll(x => x.BasePanel.FixedPanelID == fixedpanelid);
            Home.Instance.UpdatePanelData();
        }

        [JSInvokable]
        public static void OpenPanel(string panelid) {
            var panel = Home.Instance.AddPanel(panelid);
            App.SavePanelData(Home.panels);
            Home.Instance.UpdatePanelData();
        }

        [JSInvokable]
        public static string GetPanelSize(string uniquepanelid)
        {
            var panel = Home.panels.FirstOrDefault(x => x.UniquePanelID == uniquepanelid);
            return panel == null ? "" : $"{panel.PanelWidth},{panel.PanelHeight}";
        }

        [JSInvokable]
        public static void SetPanelSize(string uniquepanelid, double width, double height)
        {
            var panel = Home.panels.FirstOrDefault(x => x.UniquePanelID == uniquepanelid);
            if (panel != null)
            {
                panel.PanelWidth = width;
                panel.PanelHeight = height;
                Home.Instance.UpdatePanelData();
            }
            App.SavePanelData(Home.panels);
        }

        [JSInvokable]
        public static string GetPanelPosition(string uniquepanelid)
        {
            var panel = Home.panels.FirstOrDefault(x => x.UniquePanelID == uniquepanelid);
            return panel == null ? "" : $"{panel.PanelX},{panel.PanelY}";
        }

        [JSInvokable]
        public static void SetPanelPosition(string uniquepanelid, double x, double y)
        {
            var panel = Home.panels.FirstOrDefault(x => x.UniquePanelID == uniquepanelid);
            if (panel != null)
            {
                panel.PanelX = x;
                panel.PanelY = y;
                Home.Instance.UpdatePanelData();
            }
            App.SavePanelData(Home.panels);
        }

        [JSInvokable]
        public static string GetPanelVisibility(string uniquepanelid)
        {
            var panel = Home.panels.FirstOrDefault(x => x.UniquePanelID == uniquepanelid);
            return panel == null ? "" : panel.Deleted ? "hidden" : "visible";
        }

        [JSInvokable]
        public static void SetPanelData(string uniquepanelid, string data)
        {
            var panel = Home.panels.FirstOrDefault(x => x.UniquePanelID == uniquepanelid);
            if (panel != null)
            {
                panel.UniquePanelID = data;
                Home.Instance.UpdatePanelData();
            }
        }

        [JSInvokable]
        public static string GetPanelData(string uniquepanelid)
        {
            var panel = Home.panels.FirstOrDefault(x => x.UniquePanelID == uniquepanelid);
            return panel == null ? "" : panel.UniquePanelID;
        }

        [JSInvokable]
        public static void SetPanelHeader(string uniquepanelid, bool titlebarvisible, bool enableresize, bool enableclose, bool enabledrag)
        {
            var panel = Home.panels.FirstOrDefault(x => x.UniquePanelID == uniquepanelid);
            var exists = Home.panelHeaders.Any(x => x.persistentpanelid == uniquepanelid);
            if (exists)
            {
                var res = Home.panelHeaders.FirstOrDefault(x => x.persistentpanelid==uniquepanelid)!;
                res.titlebarvisible = titlebarvisible;
                res.enableresize = enableresize;
                res.enableclose = enableclose;
                res.enabledrag = enabledrag;
            } else
            {
                var header = new PanelHeader();
                header.persistentpanelid = uniquepanelid;
                header.titlebarvisible = titlebarvisible;
                header.enableresize = enableresize;
                header.enableclose = enableclose;
                header.enabledrag = enabledrag;
                Home.panelHeaders.Add(header);
            }
        }

        [JSInvokable]
        public static string GetPanelHeader(string persistentpanelid)
        {
            var header = Home.panelHeaders.FirstOrDefault(x => x.persistentpanelid == persistentpanelid);
            return header == null ? "" : $"{header.titlebarvisible},{header.enableresize},{header.enableclose},{header.enabledrag}";
        }


    }
}
