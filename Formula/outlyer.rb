class Outlyer < Formula
  desc "Outlyer CLI allows to easily manage your Outlyer account via command line."
  homepage "https://www.outlyer.com/"
  url "https://github.com/outlyerapp/outlyer-cli/releases/download/0.1.0/outlyer_0.1.0_Darwin_x86_64.tar.gz"
  version "0.1.0"
  sha256 "519bd7271c53c05abb86cb21058c29890ca7cde495df0eb8c830473a17121fd3"

  def install
    bin.install "outlyer"
  end
end
