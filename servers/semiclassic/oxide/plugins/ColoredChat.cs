using Oxide.Core;
using Oxide.Core.Libraries;
using System;
using System.Collections.Generic;
using System.Linq;
namespace Oxide.Plugins
{
    [Info("ColoredChat", "TexHik", "0.1")]
    [Description("Chat msgs&prefix colors")]
    public sealed class ColoredChat : RustPlugin
    {
        #region Config
        private static UserMsgColor DefaultMsgColor = new UserMsgColor
        {
            NameColor = "#FF8B55",
            TextColor = "#ffffff",
            Prefix = "ИГРОК",
            PrefixColor = "#42aaff"
        };
        private ChatConfig config;
        private class UserMsgColor
        {
            public string NameColor { get; set; }

            public string Prefix { get; set; }
            public string PrefixColor { get; set; }

            public string TextColor { get; set; }
        }
        private class ChatConfig
        {
            public Dictionary<string, UserMsgColor> colors = new Dictionary<string, UserMsgColor>()
            {
                {
                    "default",
                    DefaultMsgColor
                },
                {
                    "vip",
                    new UserMsgColor
                    {
                        NameColor = "#4CD4B0",
                        TextColor = "#ffffff",
                        Prefix = "VIP",
                        PrefixColor = "#3c78d8"
                    }
                },
                {
                    "gold",
                    new UserMsgColor
                    {
                        NameColor = "#F24D16",
                        TextColor = "#ffffff",
                        Prefix = "GOLD",
                        PrefixColor = "#6aa84f"
                    }
                },
                {
                    "grand",
                    new UserMsgColor
                    {
                        NameColor = "#44BBFF",
                        TextColor = "#ffffff",
                        Prefix = "GRAND",
                        PrefixColor = "#9900ff"
                    }
                },
                {
                    "astro",
                    new UserMsgColor
                    {
                        NameColor = "#00D717",
                        TextColor = "#ffffff",
                        Prefix = "ASTRO",
                        PrefixColor = "#980000"
                    }
                }
            };
        }
        #endregion Config
        string[] GetGroups(BasePlayer ply)
        {
            Permission permission = Interface.GetMod().GetLibrary<Permission>();
            return permission.GetUserGroups(ply.UserIDString);
        }
        protected override void LoadDefaultConfig()
        {
            config = new ChatConfig();
        }
        protected override void LoadConfig()
        {
            base.LoadConfig();
        }
        #region Hooks

        protected override void SaveConfig()
        {
            Config.WriteObject(config);
        }
        void OnServerInitialized()
        {
            LoadConfig();
            config = Config.ReadObject<ChatConfig>();


        }
        object OnPlayerChat(BasePlayer sender, string message, ConVar.Chat.ChatChannel channel)
        {
            string[] groups = GetGroups(sender);
            UserMsgColor color;
            try
            {
                color = config.colors[config.colors.Keys.ToList().Last(p => groups.Contains(p))];
            }
            catch (Exception)
            {
                color = DefaultMsgColor;
            }
            string name;
            if (string.IsNullOrEmpty(color.Prefix) || string.IsNullOrEmpty(color.PrefixColor))
            {
                name = string.Format("<color={0}>{1}</color>", color.NameColor, sender.displayName);
            } else
            {
                name = string.Format("<color={0}>[{1}]</color> <color={2}>{3}</color>", color.PrefixColor, color.Prefix, color.NameColor, sender.displayName);
            }

            string msg = string.Format("<color={0}>{1}</color>", color.TextColor, message);
			foreach (var ply in BasePlayer.activePlayerList)
            {

                ply?.SendConsoleCommand("chat.add 0", sender.userID, string.IsNullOrEmpty(sender.displayName) ? $"{msg}" : $"{name}: {msg}");
            }
            return false;
        }
        #endregion Hooks
    }
}
